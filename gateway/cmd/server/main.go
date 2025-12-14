package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"gateway/internal/config"
	"gateway/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil { log.Fatal(err) }

    rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
    r := chi.NewRouter()

    // GLOBAL MIDDLEWARE: Header Sanitization
    // This runs on EVERY request to ensure no attacker can inject trusted headers.
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Strip headers that are meant to be internal-only
            r.Header.Del("X-User-ID")
            next.ServeHTTP(w, r)
        })
    })

    r.Use(middleware.RateLimiter(rdb, cfg.RateLimitReq, cfg.RateLimitWindow))

    // Public Routes
    r.Mount("/auth/login", proxy(viper.GetString("AUTH_SERVICE_URL") + "/login"))
    r.Mount("/auth/register", proxy(viper.GetString("AUTH_SERVICE_URL") + "/register"))
    r.Mount("/events", proxy(viper.GetString("EVENT_SERVICE_URL")))
    r.Mount("/events", proxy(cfg.EventServiceURL))
    
    // Protected Routes (Apply Auth Middleware)
    r.Group(func(protected chi.Router) {
        protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
        
        // Route to Order Service (Flash Sale)
        protected.Mount("/orders", proxy(cfg.OrderServiceURL))
        protected.Mount("/book", proxy(cfg.OrderServiceURL)) // Direct mapping for /book if needed
        protected.Mount("/legacy", proxy(cfg.MonolithURL))
        protected.Mount("/users", proxy(cfg.MonolithURL)) 
    })

    log.Printf("Gateway running on :%s", cfg.Port)
    http.ListenAndServe(":"+cfg.Port, r)
}

func proxy(target string) http.HandlerFunc {
    url, _ := url.Parse(target)
    proxy := httputil.NewSingleHostReverseProxy(url)
    return func(w http.ResponseWriter, r *http.Request) {
        r.URL.Host = url.Host
        r.URL.Scheme = url.Scheme
        r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
        proxy.ServeHTTP(w, r)
    }
}