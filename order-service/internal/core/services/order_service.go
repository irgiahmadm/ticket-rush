package services

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/core/domain"
	"order-service/internal/core/ports"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type OrderService struct { 
	repo ports.OrderRepository; 
	redis *redis.Client; 
	bus ports.EventBus 
}

func NewOrderService(r ports.OrderRepository, rd *redis.Client, b ports.EventBus) *OrderService {
    return &OrderService{repo: r, redis: rd, bus: b}
}

func (s *OrderService) BookSeat(ctx context.Context, uid uuid.UUID, eid int, sid string) error {
    key := fmt.Sprintf("event:%d:available_seats", eid)
		if res, _ := s.redis.SRem(ctx, key, sid).Result(); 
		res == 0 { return errors.New("seat taken") }

    order := domain.Order{UserID: uid, EventID: eid, SeatID: sid, Status: "Pending"}
    oid, err := s.repo.Save(ctx, order)

    if err != nil { 
			s.redis.SAdd(ctx, key, sid); 
			return err 
		}

    s.bus.PublishOrderCreated(ctx, oid, uid, 100.00)
    return nil
}
