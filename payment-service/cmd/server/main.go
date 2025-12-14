package main

import (
	"context"
	"encoding/json"
	"log"
	"payment-service/internal/adapters/repository"
	"payment-service/internal/core/services"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env"); 
	viper.AutomaticEnv(); 
	viper.ReadInConfig()

	db, err := pgxpool.New(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil { log.Fatal("Failed to connect to DB:", err) }

	repo := repository.NewPostgresPaymentRepository(db)
	svc := services.NewPaymentService(repo)

	// RabbitMQ Connection with Retry
	var conn *amqp.Connection
	for i := 0; i < 15; i++ {
			conn, err = amqp.Dial(viper.GetString("RABBITMQ_URL"))
			if err == nil {
					break
			}
			log.Printf("Failed to connect to RabbitMQ (attempt %d/15): %v. Retrying in 2s...", i+1, err)
			time.Sleep(2 * time.Second)
	}
	if err != nil { log.Fatal("Failed to connect to RabbitMQ after retries:", err) }

	ch, _ := conn.Channel()
	q, _ := ch.QueueDeclare("order_created", true, false, false, false, nil)
	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	log.Println("Payment Worker Started")
	for d := range msgs {
			var req struct { OrderID int `json:"order_id"`; Amount float64 `json:"amount"` }
			json.Unmarshal(d.Body, &req)
			svc.Process(context.Background(), req.OrderID, req.Amount)
	}
}
