package postgres

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"strings"
	"workout/internal/dto"
	"workout/internal/entity"
)

func (p *postgres) CreateWorkout(ctx context.Context, workout *entity.Workout) (*entity.Workout, error) {
	query := `
        INSERT INTO workouts (
            user_id, name, sport_type, date, duration,
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

	if err := row.Scan(&workout.ID); err != nil {
		return nil, err
	}

	return workout, nil
}

func (p *postgres) GetWorkouts(ctx context.Context, userID uuid.UUID) ([]entity.Workout, error) {
	query := `
		SELECT id, user_id, name, sport_type, date, duration, distance,
    		avg_pace, avg_heart_rate, max_heart_rate, avg_cadence,
    		calories, description, created_at, updated_at
		FROM workouts
		WHERE user_id = $1
		ORDER BY date DESC;`

	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []entity.Workout
	for rows.Next() {
		var w entity.Workout
		if err := rows.Scan(
			&w.ID, &w.UserID, &w.Name, &w.SportType, &w.Date, &w.Duration,
			&w.Distance, &w.AvgPace, &w.AvgHeartRate, &w.MaxHeartRate,
			&w.AvgCadence, &w.Calories, &w.Description, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, rows.Err()
}

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

func (p *postgres) UpdateWorkout(ctx context.Context, u dto.UpdateWorkout) error {
	nextIdx := 2                        // счётчик плейсхолдеров $2, $3 …
	setClauses := make([]string, 0, 12) // сюда будем складывать "col = $n"
	args := []interface{}{u.ID}

	// Хэлпер add(…) прячет всю «механику»:
	//  • кладём "col = $n" в слайс setClauses
	//  • добавляем само значение в args
	//  • увеличиваем nextIdx
	add := func(col string, val interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, nextIdx))
		args = append(args, val)
		nextIdx++
	}

	if u.Name != nil {
		add("name", *u.Name)
	}
	if u.SportType != nil {
		add("type", *u.SportType)
	}
	if u.Duration != nil {
		add("duration", *u.Duration)
	}
	if u.Distance != nil {
		add("distance", *u.Distance)
	}
	if u.AvgPace != nil {
		add("avg_pace", *u.AvgPace)
	}
	if u.AvgHeartRate != nil {
		add("avg_heart_rate", *u.AvgHeartRate)
	}
	if u.MaxHeartRate != nil {
		add("max_heart_rate", *u.MaxHeartRate)
	}
	if u.AvgCadence != nil {
		add("avg_cadence", *u.AvgCadence)
	}
	if u.Calories != nil {
		add("calories", *u.Calories)
	}
	if u.Description != nil {
		add("description", *u.Description)
	}

	if len(setClauses) == 0 {
		return nil
	}

	// Склеиваем окончательный SQL
	// Пример получится такой:
	// UPDATE workout
	// SET name = $2, distance = $3, updated_at = CURRENT_TIMESTAMP
	// WHERE id = $1;
	query := fmt.Sprintf(`
        UPDATE workout
        SET %s,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1`, strings.Join(setClauses, ", "))

	_, err := p.db.Exec(ctx, query, args...)
	return err
}
