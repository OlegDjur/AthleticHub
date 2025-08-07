package activity

import (
	"context"
	"workout/internal/entity"
)

type Activity interface {
	CreateWorkout(ctx context.Context, workout entity.Workout) error
	GetWorkoutByID(ctx context.Context, id int64) (*entity.Workout, error)
}
