package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/room"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound         = errors.New("booking not found")
	ErrRoomNotFound     = errors.New("room not found")
	ErrRoomUnavailable  = errors.New("room is not available for the selected dates")
	ErrAlreadyCancelled = errors.New("booking is already cancelled")
)

type CreateInput struct {
	RoomID   uuid.UUID
	UserID   uuid.UUID
	CheckIn  time.Time
	CheckOut time.Time
}

// Service manages booking operations. It holds a *pgxpool.Pool rather than
// a DBTX because it is responsible for beginning transactions.
type Service struct {
	db          *pgxpool.Pool
	bookingRepo *Repository
	roomRepo    *room.Repository
}

func NewService(db *pgxpool.Pool, bookingRepo *Repository, roomRepo *room.Repository) *Service {
	return &Service{
		db:          db,
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
	}
}

// Create atomically checks room availability and creates a confirmed booking.
// A row-level lock on the room serializes concurrent requests, preventing
// overbooking when multiple requests target the same room and dates.
func (s *Service) Create(ctx context.Context, input CreateInput) (Booking, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Booking{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // no-op if Commit succeeds

	txRoomRepo := room.NewRepository(tx)
	txBookingRepo := NewRepository(tx)

	r, err := txRoomRepo.GetByIDForUpdate(ctx, input.RoomID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Booking{}, ErrRoomNotFound
	}
	if err != nil {
		return Booking{}, fmt.Errorf("lock room: %w", err)
	}

	count, err := txBookingRepo.GetOverlappingCount(ctx, input.RoomID, input.CheckIn, input.CheckOut)
	if err != nil {
		return Booking{}, fmt.Errorf("check availability: %w", err)
	}

	if count >= r.Quantity {
		return Booking{}, ErrRoomUnavailable
	}

	b, err := txBookingRepo.Create(ctx, Booking{
		RoomID:   input.RoomID,
		UserID:   input.UserID,
		Status:   "confirmed",
		CheckIn:  input.CheckIn,
		CheckOut: input.CheckOut,
	})
	if err != nil {
		return Booking{}, fmt.Errorf("create booking: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Booking{}, fmt.Errorf("commit: %w", err)
	}

	return b, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (Booking, error) {
	b, err := s.bookingRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Booking{}, ErrNotFound
	}
	if err != nil {
		return Booking{}, fmt.Errorf("get booking: %w", err)
	}

	return b, nil
}

// Cancel transitions a confirmed booking to cancelled.
// A transaction with a read-then-write prevents duplicate cancellation
// under concurrent requests targeting the same booking.
func (s *Service) Cancel(ctx context.Context, id uuid.UUID) (Booking, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Booking{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // no-op if Commit succeeds

	txBookingRepo := NewRepository(tx)

	b, err := txBookingRepo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Booking{}, ErrNotFound
	}
	if err != nil {
		return Booking{}, fmt.Errorf("get booking: %w", err)
	}

	if b.Status == "cancelled" {
		return Booking{}, ErrAlreadyCancelled
	}

	b, err = txBookingRepo.UpdateStatus(ctx, id, "cancelled")
	if err != nil {
		return Booking{}, fmt.Errorf("cancel booking: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Booking{}, fmt.Errorf("commit: %w", err)
	}

	return b, nil
}
