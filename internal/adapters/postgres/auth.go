package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"workout/internal/entity"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

func (s *postgres) CreateUser(ctx context.Context, email string, passHash []byte) (string, error) {
	const op = "storage.sqlite.SaveUser"

	var id uuid.UUID
	err := s.db.
		QueryRow(ctx, `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`, email, passHash).
		Scan(&id)
	if err != nil {
		// todo: переделать на ON CONFLICT
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}

func (s *postgres) GetUser(ctx context.Context, email string) (*entity.User, error) {
	const op = "storage.sqlite.User"

	var u entity.User
	err := s.db.
		QueryRow(ctx, `SELECT id, email, password_hash FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &u, nil
}

func (s *postgres) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	var isAdmin bool
	err := s.db.
		QueryRow(ctx, `SELECT is_admin FROM users WHERE id = $1`, userID).
		Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
