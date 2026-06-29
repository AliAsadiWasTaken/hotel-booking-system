package room

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID            uuid.UUID `json:"id"`
	HotelID       uuid.UUID `json:"hotel_id"`
	Name          string    `json:"name"`
	BedCount      int       `json:"bed_count"`
	Capacity      int       `json:"capacity"`
	Quantity      int       `json:"quantity"`
	PricePerNight float64   `json:"price_per_night"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
