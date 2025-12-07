package repository

import (
	"context"
	"order-service/internal/core/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)	

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, order *domain.Order) (int, error) {
	// PGX Transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	// Defer rollback (safe if commit isn't reached)
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx, 
		"INSERT INTO orders (user_id, product_id, status) VALUES ($1, $2, $3) RETURNING id", 
		order.UserID, order.ProductID, order.Status).Scan(&id)
	
	if err != nil {
		return 0, err
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}
