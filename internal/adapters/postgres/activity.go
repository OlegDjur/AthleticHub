package postgres

import (
	"context"
	"workout/internal/entity"
)

func (p *postgres) CreateWorkout(ctx context.Context, workout entity.Workout) error {
	query := `
        INSERT INTO workouts (
            user_id, name, type, start_time, duration,
            distance, avg_pace, avg_heart_rate, max_heart_rate,
            avg_cadence, calories, description, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        ) RETURNING id`

	row := p.db.QueryRow(ctx, query,
		workout.UserID, workout.Name, workout.SportType, workout.Date,
		workout.Duration, workout.Distance,
		workout.AvgPace, workout.AvgHeartRate, workout.MaxHeartRate,
		workout.AvgCadence, workout.Calories, workout.Description,
		workout.CreatedAt, workout.UpdatedAt,
	)
	return row.Scan(&workout.ID)
}

// func (r *WorkoutRepository) GetWorkouts(ctx context.Context, id int64) (*models.Workout, error) {
// 	return nil, nil
// }

func (p *postgres) GetWorkoutByID(ctx context.Context, id int64) (*entity.Workout, error) {
	query := `
        SELECT id, name, type, date, duration, duration,
            distance, avg_pace,  avg_heart_rate, max_heart_rate,
            avg_power, max_power, avg_cadence, calories,
            description,  created_at, updated_at
        FROM workouts WHERE id = $1`

	workout := &entity.Workout{}
	err := p.db.QueryRow(ctx, query, id).Scan(
		&workout.Name, &workout.SportType, &workout.Date,
		&workout.Duration, &workout.Duration, &workout.Distance,
		&workout.AvgPace, &workout.AvgHeartRate, &workout.MaxHeartRate,
		&workout.AvgCadence, &workout.Calories, &workout.Description,
		&workout.CreatedAt, &workout.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (p *postgres) GetWorkoutsByUserID(ctx context.Context, userID int64) ([]*entity.Workout, error) {
	query := `
        SELECT id, user_id, name, type, duration, distance, avg_heart_rate, calories, created_at
        FROM workouts 
        WHERE user_id = $1 
        ORDER BY start_time DESC`

	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []*entity.Workout
	for rows.Next() {
		workout := &entity.Workout{}
		err := rows.Scan(
			&workout.ID, &workout.UserID, &workout.Name, &workout.SportType,
			&workout.Duration, &workout.Distance,
			&workout.AvgHeartRate, &workout.Calories, &workout.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}
	return workouts, nil
}
