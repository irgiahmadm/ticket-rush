package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Standard JSON response struct for Gateway Middleware
type ErrorResponse struct {
    StatusCode int         `json:"statusCode"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data"` // No omitempty
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(ErrorResponse{StatusCode: 401, Message: "Unauthorized: No token header", Data: nil})
                return
            }
            
            tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
            token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })
            
            if err != nil || !token.Valid {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(ErrorResponse{StatusCode: 401, Message: "Unauthorized: Invalid token", Data: nil})
                return
            }

            // STRATEGY IMPLEMENTATION: Identity Propagation
            // Extract user_id from claims
            if claims, ok := token.Claims.(jwt.MapClaims); ok {
                if userID, ok := claims["user_id"].(string); ok {
                    // Inject into header for downstream services
                    r.Header.Set("X-User-ID", userID)
                }
            }

            next.ServeHTTP(w, r)
        })
    }
}