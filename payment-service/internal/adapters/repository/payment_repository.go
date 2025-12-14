package repository

import (
	"context"
	"payment-service/internal/core/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPaymentRepository struct{ pool *pgxpool.Pool }

func NewPostgresPaymentRepository(p *pgxpool.Pool) ports.PaymentRepository {
	return &PostgresPaymentRepository{pool: p}
}

func (r *PostgresPaymentRepository) SavePayment(ctx context.Context, oid int, status string) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO payments (order_id, amount, status) VALUES ($1, 100.00, $2)", oid, status)
	return err
}