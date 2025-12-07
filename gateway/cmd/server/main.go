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
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil { log.Fatal(err) }

    rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
    r := chi.NewRouter()
    r.Use(middleware.RateLimiter(rdb, cfg.RateLimitReq, cfg.RateLimitWindow))

    r.Mount("/auth", http.StripPrefix("/auth", proxy(cfg.AuthServiceURL)))

    r.Group(func(p chi.Router) {
        p.Use(middleware.Auth(cfg.JWTSecret))

        p.Mount("/orders", proxy(cfg.OrderServiceURL))
        p.Mount("/legacy", proxy(cfg.MonolithURL))
        p.Mount("/users", proxy(cfg.MonolithURL)) 
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
