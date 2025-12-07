package token

import (
	"auth-service/internal/core/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secret string
}

func NewJWTGenerator(secret string) *JWTGenerator {
	return &JWTGenerator{secret: secret}
}

func (j *JWTGenerator) GenerateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(j.secret))
}
