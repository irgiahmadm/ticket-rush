package repository

import (
	"context"
	"event-service/internal/core/domain"
	"event-service/internal/core/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresEventRepository struct{ pool *pgxpool.Pool }

func NewPostgresEventRepository(p *pgxpool.Pool) ports.EventRepository {
	return &PostgresEventRepository{pool: p}
}

func (r *PostgresEventRepository) FindAll() ([]domain.Event, error) {
	rows, _ := r.pool.Query(context.Background(), "SELECT id, name, date FROM events")
	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		rows.Scan(&e.ID, &e.Name, &e.Date)
		events = append(events, e)
	}
	return events, nil
}