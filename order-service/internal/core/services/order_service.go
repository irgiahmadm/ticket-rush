package services

import (
	"context"
	"order-service/internal/core/domain"
	"order-service/internal/core/ports"
)

type Service struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PlaceOrder(ctx context.Context, userID string, productID int) (*domain.Order, error) {
	order := &domain.Order{
		UserID:    userID,
		ProductID: productID,
		Status:    "PENDING",
	}
	id, err := s.repo.Save(ctx, order)
	if err != nil {
		return nil, err
	}
	order.ID = id
	return order, nil
}
