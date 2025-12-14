package messaging

import (
	"context"
	"encoding/json"
	"log"
	"order-service/internal/core/ports"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBus struct { ch *amqp.Channel }

func NewRabbitMQBus(conn *amqp.Connection) ports.EventBus {
    ch, err := conn.Channel()
    if err != nil {
        log.Fatal("Failed to open channel:", err)
    }
    ch.QueueDeclare("order_created", true, false, false, false, nil)
    return &RabbitMQBus{ch: ch}
}
func (rabbitMq *RabbitMQBus) PublishOrderCreated(ctx context.Context, orderID int, uid uuid.UUID, amount float64) error {
    body, _ := json.Marshal(map[string]interface{}{"order_id": orderID, "user_id": uid, "amount": amount})
    return rabbitMq.ch.PublishWithContext(ctx, "", "order_created", false, false, amqp.Publishing{ContentType: "application/json", Body: body})
}