package services

import (
	"auth-service/internal/core/domain"
	"auth-service/internal/core/ports"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo ports.UserRepository
	tokenGen ports.TokenGenerator
}

func NewAuthService(repo ports.UserRepository, tokenGen ports.TokenGenerator) *Service {
	return &Service{repo: repo, tokenGen: tokenGen}
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Fetch User (Password in DB is hashed)
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. Verify Password using Bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Generate Token
	return s.tokenGen.GenerateToken(user)
}

func (s *Service) Register(ctx context.Context, email, password string) error {
	// 1. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.repo.Save(ctx, user)
}