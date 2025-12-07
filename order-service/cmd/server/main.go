package main

import (
	"context"
	"log"

	"order-service/internal/adapters/handler"
	"order-service/internal/adapters/repository"
	"order-service/internal/config"
	"order-service/internal/core/services"
	"order-service/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	// 1. Infrastructure (PGX Pool)
	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer dbPool.Close()

	// Dependency Injection (Hexagonal Wiring)
	repo := repository.NewPostgresRepo(dbPool)
	svc := services.NewService(repo)
	h := handler.NewHandler(svc)

	r := gin.Default()
	r.POST("/orders", response.Wrap(h.CreateOrder))

	r.Run(":" + cfg.Port)
}