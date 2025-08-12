package controller

import (
	"context"
	"workout/internal/dto"
)

type Activity interface {
	UploadFile(ctx context.Context, data []byte) (*dto.UploadFile, error)
	CreateWorkout(ctx context.Context, w dto.WorkoutDTO) error
	UpdateWorkout(ctx context.Context, u *dto.UpdateWorkout) error
}

type Auth interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RegisterNewUser(ctx context.Context, req dto.RegisterRequest) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
