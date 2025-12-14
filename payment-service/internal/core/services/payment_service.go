package services

import (
	"context"
	"log"
	"payment-service/internal/core/ports"
	"time"
)

type PaymentService struct{ repo ports.PaymentRepository }

func NewPaymentService(r ports.PaymentRepository) *PaymentService { return &PaymentService{repo: r} }

func (s *PaymentService) Process(ctx context.Context, orderID int, amount float64) {
	log.Printf("Processing payment for Order %d", orderID)
	time.Sleep(time.Second) // Simulate bank
	s.repo.SavePayment(ctx, orderID, "Success")
}