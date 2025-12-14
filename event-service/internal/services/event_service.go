package services

import (
	"event-service/internal/core/domain"
	"event-service/internal/core/ports"
)

type EventService struct{ repo ports.EventRepository }

func NewEventService(r ports.EventRepository) *EventService { return &EventService{repo: r} }
func (s *EventService) GetAll() ([]domain.Event, error)     { return s.repo.FindAll() }
