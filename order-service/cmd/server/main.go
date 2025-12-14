package main

import (
	"context"
	"log"

	"order-service/internal/adapters/handler"
	"order-service/internal/adapters/messaging"
	"order-service/internal/adapters/repository"
	"order-service/internal/core/services"
	"order-service/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	viper.SetConfigFile(".env"); viper.AutomaticEnv(); viper.ReadInConfig()
    db, err := pgxpool.New(context.Background(), viper.GetString("DATABASE_URL"))
    if err != nil { log.Fatal("Failed to connect to DB:", err) }

    rdb := redis.NewClient(&redis.Options{Addr: viper.GetString("REDIS_ADDR")})
    conn, err := amqp.Dial(viper.GetString("RABBITMQ_URL"))
    if err != nil { log.Fatal("Failed to connect to RabbitMQ:", err) }

    repo := repository.NewPostgresOrderRepository(db)
    bus := messaging.NewRabbitMQBus(conn)
    svc := services.NewOrderService(repo, rdb, bus)
    h := handler.NewOrderHandler(svc)

    r := gin.Default()
    r.POST("/book", response.Wrap(h.BookSeat))
    r.Run(":" + viper.GetString("PORT"))
}