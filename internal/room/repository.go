package room

import (
	"context"
	"fmt"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/database"
	"github.com/google/uuid"
)

type Repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, room Room) (Room, error) {
	query := `
		INSERT INTO rooms (hotel_id, name, bed_count, capacity, quantity, price_per_night)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, hotel_id, name, bed_count, capacity, quantity, price_per_night, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, room.HotelID, room.Name, room.BedCount, room.Capacity, room.Quantity, room.PricePerNight).
		Scan(&room.ID, &room.HotelID, &room.Name, &room.BedCount, &room.Capacity, &room.Quantity, &room.PricePerNight, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return Room{}, fmt.Errorf("create room: %w", err)
	}

	return room, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (Room, error) {
	var room Room

	query := `
		SELECT id, hotel_id, name, bed_count, capacity, quantity, price_per_night, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&room.ID, &room.HotelID, &room.Name, &room.BedCount, &room.Capacity, &room.Quantity, &room.PricePerNight, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return Room{}, fmt.Errorf("get room by id: %w", err)
	}

	return room, nil
}

// GetByIDForUpdate locks the room row for the duration of the calling transaction.
// Use this when you need to check availability and create a booking atomically —
// it prevents another transaction from reading or modifying this row until the
// current transaction commits or rolls back.
func (r *Repository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (Room, error) {
	var room Room

	query := `
		SELECT id, hotel_id, name, bed_count, capacity, quantity, price_per_night, created_at, updated_at
		FROM rooms
		WHERE id = $1
		FOR UPDATE
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&room.ID, &room.HotelID, &room.Name, &room.BedCount, &room.Capacity, &room.Quantity, &room.PricePerNight, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return Room{}, fmt.Errorf("get room by id for update: %w", err)
	}

	return room, nil
}

func (r *Repository) ListByHotelID(ctx context.Context, hotelID uuid.UUID) ([]Room, error) {
	query := `
		SELECT id, hotel_id, name, bed_count, capacity, quantity, price_per_night, created_at, updated_at
		FROM rooms
		WHERE hotel_id = $1
	`

	rows, err := r.db.Query(ctx, query, hotelID)
	if err != nil {
		return nil, fmt.Errorf("list rooms by hotel id: %w", err)
	}
	defer rows.Close()

	var rooms []Room

	for rows.Next() {
		var room Room
		err := rows.Scan(&room.ID, &room.HotelID, &room.Name, &room.BedCount, &room.Capacity, &room.Quantity, &room.PricePerNight, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan room: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list rooms by hotel id: %w", err)
	}

	return rooms, nil
}
