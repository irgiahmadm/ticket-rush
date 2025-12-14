package domain

import "github.com/google/uuid"

type Order struct {
	ID      int
	UserID  uuid.UUID
	EventID int
	SeatID  string
	Status  string
}
