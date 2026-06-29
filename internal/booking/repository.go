package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, booking Booking) (Booking, error) {
	query := `
		INSERT INTO bookings (room_id, user_id, status, check_in, check_out)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, room_id, user_id, status, check_in, check_out, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, booking.RoomID, booking.UserID, booking.Status, booking.CheckIn, booking.CheckOut).
		Scan(&booking.ID, &booking.RoomID, &booking.UserID, &booking.Status, &booking.CheckIn, &booking.CheckOut, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		return Booking{}, fmt.Errorf("create booking: %w", err)
	}

	return booking, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (Booking, error) {
	var booking Booking

	query := `
		SELECT id, room_id, user_id, status, check_in, check_out, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&booking.ID, &booking.RoomID, &booking.UserID, &booking.Status, &booking.CheckIn, &booking.CheckOut, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		return Booking{}, fmt.Errorf("get booking by id: %w", err)
	}

	return booking, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (Booking, error) {
	var booking Booking

	query := `
		UPDATE bookings
		SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, room_id, user_id, status, check_in, check_out, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, status, id).
		Scan(&booking.ID, &booking.RoomID, &booking.UserID, &booking.Status, &booking.CheckIn, &booking.CheckOut, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		return Booking{}, fmt.Errorf("update booking status: %w", err)
	}

	return booking, nil
}

func (r *Repository) GetOverlappingCount(ctx context.Context, roomID uuid.UUID, checkIn, checkOut time.Time) (int, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM bookings
		WHERE room_id = $1
		AND status = 'confirmed'
		AND check_in  < $2
		AND check_out > $3
	`

	err := r.db.QueryRow(ctx, query, roomID, checkOut, checkIn).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get overlapping count: %w", err)
	}

	return count, nil
}
