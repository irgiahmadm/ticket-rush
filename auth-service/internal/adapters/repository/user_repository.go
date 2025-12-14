package repository

import (
	"auth-service/internal/core/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = $1"

	var user domain.User
	// pgx uses context for queries
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user domain.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2)"
	_, err := r.db.Exec(ctx, query, user.Email, user.Password)
	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
    var u domain.User
    err := r.db.QueryRow(ctx, "SELECT id, email, password FROM users WHERE id=$1", id).Scan(&u.ID, &u.Email, &u.Password)
    return u, err
}
