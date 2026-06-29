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

// Service holds the raw *pgxpool.Pool — not a DBTX.
// It needs the pool specifically to begin transactions.
// Repositories only need DBTX, but the service is the one
// that decides when a transaction is needed.
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

// Create books a room atomically using a transaction.
//
// The flow inside the transaction:
//  1. Lock the room row with FOR UPDATE — prevents concurrent bookings
//     from seeing stale availability until this transaction commits.
//  2. Count overlapping confirmed bookings.
//  3. Compare against room quantity — reject if fully booked.
//  4. Insert the booking.
//  5. Commit — all four steps succeed or none of them do.
//
// Without the transaction, two requests could both pass the availability
// check and both create a booking, causing overbooking.
func (s *Service) Create(ctx context.Context, input CreateInput) (Booking, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Booking{}, fmt.Errorf("begin transaction: %w", err)
	}
	// defer Rollback is always safe — if Commit already succeeded,
	// pgx treats a subsequent Rollback as a no-op.
	defer tx.Rollback(ctx)

	// Temporary repositories scoped to this transaction.
	// They use tx (which satisfies DBTX) instead of the pool.
	// The same repository code runs — it has no idea it's in a transaction.
	txRoomRepo := room.NewRepository(tx)
	txBookingRepo := NewRepository(tx)

	// Step 1 — lock the room row.
	// Any other transaction trying to lock this row must wait until we commit.
	r, err := txRoomRepo.GetByIDForUpdate(ctx, input.RoomID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Booking{}, ErrRoomNotFound
	}
	if err != nil {
		return Booking{}, fmt.Errorf("lock room: %w", err)
	}

	// Step 2 — count overlapping confirmed bookings.
	// This read is safe because we hold the lock on the room row.
	count, err := txBookingRepo.GetOverlappingCount(ctx, input.RoomID, input.CheckIn, input.CheckOut)
	if err != nil {
		return Booking{}, fmt.Errorf("check availability: %w", err)
	}

	// Step 3 — reject if no slots remain.
	if count >= r.Quantity {
		return Booking{}, ErrRoomUnavailable
	}

	// Step 4 — create the booking.
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

	// Step 5 — commit. If this fails, defer fires and rolls back.
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

// Cancel cancels a confirmed booking.
// A transaction is used here to prevent a race condition where two
// requests cancel the same booking simultaneously — both would check
// the status, both would see "confirmed", and both would try to cancel.
func (s *Service) Cancel(ctx context.Context, id uuid.UUID) (Booking, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Booking{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

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
