package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"gateway/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func main() {
    viper.SetConfigFile(".env"); viper.AutomaticEnv(); viper.ReadInConfig()
    r := chi.NewRouter()
    rdb := redis.NewClient(&redis.Options{Addr: viper.GetString("REDIS_ADDR")})

    // 1. Global Middleware
    r.Use(func(next http.Handler) http.Handler { 
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            r.Header.Del("X-User-ID") // Sanitizer
            next.ServeHTTP(w, r)
        })
    })
    r.Use(middleware.RateLimiter(rdb, 100, 60))

    // 2. Public Routes
    r.Mount("/auth/login", http.StripPrefix("/auth", proxy(viper.GetString("AUTH_SERVICE_URL"))))
    r.Mount("/auth/register", http.StripPrefix("/auth", proxy(viper.GetString("AUTH_SERVICE_URL"))))
    r.Mount("/events", proxy(viper.GetString("EVENT_SERVICE_URL")))
    
    // 3. Protected Routes
    r.Group(func(p chi.Router) {
        p.Use(middleware.AuthMiddleware(viper.GetString("JWT_SECRET")))
        p.Mount("/auth/me", http.StripPrefix("/auth", proxy(viper.GetString("AUTH_SERVICE_URL"))))
        p.Post("/book", func(w http.ResponseWriter, req *http.Request) {
            target, _ := url.Parse(viper.GetString("ORDER_SERVICE_URL"))
            proxy := httputil.NewSingleHostReverseProxy(target); proxy.ServeHTTP(w, req)
        })
    })
    log.Println("Gateway on :" + viper.GetString("PORT"))
    http.ListenAndServe(":"+viper.GetString("PORT"), r)
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