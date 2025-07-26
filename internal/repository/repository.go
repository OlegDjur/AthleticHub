package repository

import (
	"context"
	"database/sql"

	"workout/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type WorkoutRepository struct {
	db *sql.DB
}

func NewWorkoutRepository(db *sql.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

func (r *WorkoutRepository) Create(ctx context.Context, workout *models.Workout) error {
	query := `
        INSERT INTO workouts (
            user_id, name, type, sub_type, start_time, end_time, duration, moving_time,
            distance, avg_speed, max_speed, avg_pace, avg_heart_rate, max_heart_rate,
            elevation_gain, elevation_loss, min_elevation, max_elevation,
            avg_power, max_power, avg_cadence, max_cadence, calories,
            avg_temperature, min_temperature, max_temperature,
            description, device_name, file_type, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		workout.UserID, workout.Name, workout.Type, workout.SubType,
		workout.StartTime, workout.EndTime, workout.Duration, workout.MovingTime,
		workout.Distance, workout.AvgSpeed, workout.MaxSpeed, workout.AvgPace,
		workout.AvgHeartRate, workout.MaxHeartRate,
		workout.TotalElevationGain, workout.TotalElevationLoss,
		workout.MinElevation, workout.MaxElevation,
		workout.AvgPower, workout.MaxPower, workout.AvgCadence, workout.MaxCadence,
		workout.Calories, workout.AvgTemperature, workout.MinTemperature, workout.MaxTemperature,
		workout.Description, workout.DeviceName, workout.FileType,
		workout.CreatedAt, workout.UpdatedAt,
	).Scan(&workout.ID)

	return err
}

func (r *WorkoutRepository) GetByID(ctx context.Context, id int64) (*models.Workout, error) {
	query := `
        SELECT id, user_id, name, type, sub_type, start_time, end_time, duration, moving_time,
               distance, avg_speed, max_speed, avg_pace, avg_heart_rate, max_heart_rate,
               elevation_gain, elevation_loss, min_elevation, max_elevation,
               avg_power, max_power, avg_cadence, max_cadence, calories,
               avg_temperature, min_temperature, max_temperature,
               description, device_name, file_type, created_at, updated_at
        FROM workouts WHERE id = ?`

	workout := &models.Workout{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&workout.ID, &workout.UserID, &workout.Name, &workout.Type, &workout.SubType,
		&workout.StartTime, &workout.EndTime, &workout.Duration, &workout.MovingTime,
		&workout.Distance, &workout.AvgSpeed, &workout.MaxSpeed, &workout.AvgPace,
		&workout.AvgHeartRate, &workout.MaxHeartRate,
		&workout.TotalElevationGain, &workout.TotalElevationLoss,
		&workout.MinElevation, &workout.MaxElevation,
		&workout.AvgPower, &workout.MaxPower, &workout.AvgCadence, &workout.MaxCadence,
		&workout.Calories, &workout.AvgTemperature, &workout.MinTemperature, &workout.MaxTemperature,
		&workout.Description, &workout.DeviceName, &workout.FileType,
		&workout.CreatedAt, &workout.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (r *WorkoutRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Workout, error) {
	query := `
        SELECT id, user_id, name, type, start_time, duration, distance, avg_heart_rate, calories, created_at
        FROM workouts 
        WHERE user_id = ? 
        ORDER BY start_time DESC 
        LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []*models.Workout
	for rows.Next() {
		workout := &models.Workout{}
		err := rows.Scan(
			&workout.ID, &workout.UserID, &workout.Name, &workout.Type,
			&workout.StartTime, &workout.Duration, &workout.Distance,
			&workout.AvgHeartRate, &workout.Calories, &workout.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	return workouts, nil
}

// TrackPoint repository methods
func (r *WorkoutRepository) CreateTrackPoints(ctx context.Context, points []*models.TrackPoint) error {
	if len(points) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO track_points (
            workout_id, timestamp, latitude, longitude, elevation, 
            heart_rate, speed, power, cadence, temperature, distance, grade
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, point := range points {
		_, err = stmt.ExecContext(ctx,
			point.WorkoutID, point.Timestamp, point.Latitude, point.Longitude,
			point.Elevation, point.HeartRate, point.Speed, point.Power,
			point.Cadence, point.Temperature, point.Distance, point.Grade,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *WorkoutRepository) GetTrackPoints(ctx context.Context, workoutID int64) ([]*models.TrackPoint, error) {
	query := `
        SELECT id, workout_id, timestamp, latitude, longitude, elevation,
               heart_rate, speed, power, cadence, temperature, distance, grade
        FROM track_points 
        WHERE workout_id = ? 
        ORDER BY timestamp`

	rows, err := r.db.QueryContext(ctx, query, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*models.TrackPoint
	for rows.Next() {
		point := &models.TrackPoint{}
		err := rows.Scan(
			&point.ID, &point.WorkoutID, &point.Timestamp,
			&point.Latitude, &point.Longitude, &point.Elevation,
			&point.HeartRate, &point.Speed, &point.Power,
			&point.Cadence, &point.Temperature, &point.Distance, &point.Grade,
		)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, nil
}
