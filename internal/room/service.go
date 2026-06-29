package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("room not found")

type CreateInput struct {
	HotelID       uuid.UUID
	Name          string
	BedCount      int
	Capacity      int
	Quantity      int
	PricePerNight float64
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Room, error) {
	room, err := s.repo.Create(ctx, Room{
		HotelID:       input.HotelID,
		Name:          input.Name,
		BedCount:      input.BedCount,
		Capacity:      input.Capacity,
		Quantity:      input.Quantity,
		PricePerNight: input.PricePerNight,
	})
	if err != nil {
		return Room{}, fmt.Errorf("create room: %w", err)
	}

	return room, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (Room, error) {
	room, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Room{}, ErrNotFound
	}
	if err != nil {
		return Room{}, fmt.Errorf("get room: %w", err)
	}

	return room, nil
}

func (s *Service) ListByHotelID(ctx context.Context, hotelID uuid.UUID) ([]Room, error) {
	rooms, err := s.repo.ListByHotelID(ctx, hotelID)
	if err != nil {
		return nil, fmt.Errorf("list rooms: %w", err)
	}

	return rooms, nil
}
