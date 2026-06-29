package hotel

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, hotel Hotel) (Hotel, error) {
	query := `
		INSERT INTO hotels (name, address, city, country)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, address, city, country, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, hotel.Name, hotel.Address, hotel.City, hotel.Country).
		Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, &hotel.Country, &hotel.CreatedAt, &hotel.UpdatedAt)
	if err != nil {
		return Hotel{}, fmt.Errorf("create hotel: %w", err)
	}

	return hotel, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (Hotel, error) {
	var hotel Hotel

	query := `
		SELECT id, name, address, city, country, created_at, updated_at
		FROM hotels
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, &hotel.Country, &hotel.CreatedAt, &hotel.UpdatedAt)
	if err != nil {
		return Hotel{}, fmt.Errorf("get hotel by id: %w", err)
	}

	return hotel, nil
}

func (r *Repository) List(ctx context.Context) ([]Hotel, error) {
	query := `
		SELECT id, name, address, city, country, created_at, updated_at
		FROM hotels
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list hotels: %w", err)
	}
	defer rows.Close()

	var hotels []Hotel

	for rows.Next() {
		var hotel Hotel
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address, &hotel.City, &hotel.Country, &hotel.CreatedAt, &hotel.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan hotel: %w", err)
		}
		hotels = append(hotels, hotel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list hotels: %w", err)
	}

	return hotels, nil
}
