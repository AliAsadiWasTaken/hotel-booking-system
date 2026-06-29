package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound = errors.New("user not found")

	// ErrEmailTaken is returned when a registration uses an email that already exists.
	// The repository returns a Postgres unique constraint error (code 23505).
	// The service translates that into a domain error the handler understands.
	ErrEmailTaken = errors.New("email already in use")
)

type CreateInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string // plain text — hashed before storage
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (User, error) {
	// bcrypt.DefaultCost is 10 rounds — a good balance between security and speed.
	// Never store plain text passwords. Never use MD5 or SHA for passwords.
	// bcrypt is slow by design — it makes brute-force attacks expensive.
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repo.Create(ctx, User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: string(hash),
	})
	if err != nil {
		// Postgres error code 23505 = unique_violation.
		// errors.As unwraps the error chain to find a *pgconn.PgError.
		// This is how you detect specific database constraint violations.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailTaken
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
