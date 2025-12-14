package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client, limit int, window int) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := strings.Split(r.RemoteAddr, ":")[0]
            key := "rate_limit:" + ip

            ctx := r.Context()
            count, err := rdb.Incr(ctx, key).Result()
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            if count == 1 {
                rdb.Expire(ctx, key, time.Duration(window)*time.Second)
            }

            if int(count) > limit {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
