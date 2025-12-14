package ports

import (
	"context"
	"order-service/internal/core/domain"

	"github.com/google/uuid"
)

type OrderService interface { BookSeat(ctx context.Context, uid uuid.UUID, eid int, sid string) error }
type OrderRepository interface { Save(ctx context.Context, o domain.Order) (int, error) }
type EventBus interface { PublishOrderCreated(ctx context.Context, orderID int, userID uuid.UUID, amount float64) error }
