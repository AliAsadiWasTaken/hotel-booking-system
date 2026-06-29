package hotel

import (
	"time"

	"github.com/google/uuid"
)

type Hotel struct {
	ID        uuid.UUID
	Name      string
	Address   string
	City      string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
