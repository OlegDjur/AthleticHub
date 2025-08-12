package auth

import (
	"context"
	"workout/internal/entity"
)

type Auth interface {
	CreateUser(ctx context.Context, email string, passHash []byte) (string, error)
	GetUser(ctx context.Context, email string) (*entity.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
