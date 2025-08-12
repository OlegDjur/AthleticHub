package activity

import (
	"context"
	"github.com/gofrs/uuid/v5"
	"workout/internal/dto"
	"workout/internal/entity"
)

type Activity interface {
	GetWorkouts(ctx context.Context, userID uuid.UUID) ([]entity.Workout, error)
	CreateWorkout(ctx context.Context, workout *entity.Workout) (*entity.Workout, error)
	GetWorkoutByID(ctx context.Context, id int64) (*entity.Workout, error)
	UpdateWorkout(ctx context.Context, u dto.UpdateWorkout) error
}
