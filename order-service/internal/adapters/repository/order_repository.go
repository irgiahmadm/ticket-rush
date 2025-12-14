package repository

import (
	"context"
	"order-service/internal/core/domain"
	"order-service/internal/core/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)	

type PostgresOrderRepository struct { 
	pool *pgxpool.Pool 
}
func NewPostgresOrderRepository(p *pgxpool.Pool) ports.OrderRepository { 
	return &PostgresOrderRepository{pool: p} 
}
func (r *PostgresOrderRepository) Save(ctx context.Context, order domain.Order) (int, error) {
    var id int
    err := r.pool.QueryRow(ctx, "INSERT INTO orders (user_id, event_id, seat_id, amount, status) VALUES ($1, $2, $3, 100.00, $4) RETURNING id", order.UserID, order.EventID, order.SeatID, order.Status).Scan(&id)
    return id, err
}