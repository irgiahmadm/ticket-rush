package ports

import "context"

type PaymentRepository interface {
	SavePayment(ctx context.Context, orderID int, status string) error
}

type PaymentService interface {
	Process(ctx context.Context, orderID int, amount float64)
}
