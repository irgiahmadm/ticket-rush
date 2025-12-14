package ports

import "event-service/internal/core/domain"

type EventService interface {
	GetAll() ([]domain.Event, error)
}

type EventRepository interface {
	FindAll() ([]domain.Event, error)
}