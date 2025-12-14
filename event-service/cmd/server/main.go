package main

import (
	"context"
	"event-service/internal/adapters/handler"
	"event-service/internal/adapters/repository"
	"event-service/internal/response"
	"event-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
    viper.SetConfigFile(".env"); viper.AutomaticEnv(); viper.ReadInConfig()
    db, _ := pgxpool.New(context.Background(), viper.GetString("DATABASE_URL"))
    repo := repository.NewPostgresEventRepository(db)
    svc := services.NewEventService(repo)
    h := handler.NewEventHandler(svc)

    r := gin.Default()
    r.GET("/events", response.Wrap(h.GetAll))
    r.Run(":" + viper.GetString("PORT"))
}

