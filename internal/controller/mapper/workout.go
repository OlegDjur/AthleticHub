package mapper

import (
	"workout/internal/dto"
	"workout/internal/entity"
	"workout/internal/utils"
)

func ConvertWorkoutToDTO(w entity.Workout) *dto.WorkoutDTO {

	pace := utils.FormatPaceMMSS(w.AvgPace)
	return &dto.WorkoutDTO{
		ID:           w.ID,
		Name:         w.Name,
		Description:  w.Description,
		Date:         w.Date.Format("2006-01-02"),
		Duration:     w.Duration.String(),
		Distance:     w.Distance,
		AvgPace:      pace,
		AvgHeartRate: w.AvgHeartRate,
		MaxHeartRate: w.MaxHeartRate,
		AvgCadence:   w.AvgCadence,
		SportType:    w.SportType,
		Calories:     w.Calories,
	}
}
