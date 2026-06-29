package booking

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	UserID    uuid.UUID
	Status    string
	CheckIn   time.Time
	CheckOut  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
