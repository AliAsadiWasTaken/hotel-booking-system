package room

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID            uuid.UUID
	HotelID       uuid.UUID
	Name          string
	BedCount      int
	Capacity      int
	Quantity      int
	PricePerNight float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
