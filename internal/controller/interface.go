package controller

import (
	"context"
	"workout/internal/dto"
	"workout/internal/entity"
)

type Activity interface {
	UploadFile(ctx context.Context, data []byte) (*dto.UploadFile, error)
	CreateWorkout(ctx context.Context, w entity.Workout) error
	UpdateWorkout(ctx context.Context, u *dto.UpdateWorkout) error
}
