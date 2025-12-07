package ports

import (
	"context"
	"order-service/internal/core/domain"
)

type OrderService interface {
	PlaceOrder(ctx context.Context, userID string, productID int) (*domain.Order, error)
}

type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) (int, error)
}
