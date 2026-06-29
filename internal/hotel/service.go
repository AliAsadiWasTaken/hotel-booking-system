package hotel

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ErrNotFound is returned when a hotel does not exist.
// Handlers check for this to return a 404 instead of a 500.
var ErrNotFound = errors.New("hotel not found")

type CreateInput struct {
	Name    string
	Address string
	City    string
	Country string
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Hotel, error) {
	hotel, err := s.repo.Create(ctx, Hotel{
		Name:    input.Name,
		Address: input.Address,
		City:    input.City,
		Country: input.Country,
	})
	if err != nil {
		return Hotel{}, fmt.Errorf("create hotel: %w", err)
	}

	return hotel, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (Hotel, error) {
	hotel, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Hotel{}, ErrNotFound
	}
	if err != nil {
		return Hotel{}, fmt.Errorf("get hotel: %w", err)
	}

	return hotel, nil
}

func (s *Service) List(ctx context.Context) ([]Hotel, error) {
	hotels, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list hotels: %w", err)
	}

	return hotels, nil
}
