package dto

import "github.com/gofrs/uuid/v5"

type UpdateWorkout struct {
	ID           uuid.UUID
	Name         *string
	SportType    *string
	Duration     *int32
	Distance     *float32
	AvgPace      *float32
	AvgHeartRate *int16
	MaxHeartRate *int16
	AvgCadence   *int16
	Calories     *int32
	Description  *string
	DeviceName   *string
}
