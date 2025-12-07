package main

import (
	"auth-service/internal/adapters/handler"
	"auth-service/internal/adapters/repository"
	"auth-service/internal/adapters/token"
	"auth-service/internal/config"
	"auth-service/internal/core/services"
	"auth-service/internal/response"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil { log.Fatal("Config error:", err) }

    // 1. Infrastructure (PGX Pool)
    dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
    if err != nil { log.Fatal("Unable to connect to database:", err) }
    defer dbPool.Close()

    // 2. Driven Adapters (Outbound)
    repo := repository.NewPostgresRepo(dbPool)
    tokenGen := token.NewJWTGenerator(cfg.JWTSecret)

    // 3. Core Service
    svc := services.NewAuthService(repo, tokenGen)

    // 4. Driving Adapter (Inbound)
    h := handler.NewHandler(svc)

    // 5. Framework
    r := gin.Default()
    r.POST("/login", response.Wrap(h.Login))
    r.POST("/register", response.Wrap(h.Register))

    r.Run(":" + cfg.Port)
}
