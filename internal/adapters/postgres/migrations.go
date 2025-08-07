package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const CreateTablesSQL = `
-- Таблица тренировок
CREATE TABLE IF NOT EXISTS workouts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    duration INTEGER NOT NULL,
    distance DECIMAL(10,2),
    avg_pace DECIMAL(10,2),
    avg_heart_rate INTEGER,
    max_heart_rate INTEGER,
    avg_cadence INTEGER,
    calories INTEGER,
    description TEXT,
    device_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_workouts_start_time ON workouts(start_time);
CREATE INDEX IF NOT EXISTS idx_workouts_type ON workouts(type);
CREATE INDEX IF NOT EXISTS idx_track_points_workout_id ON track_points(workout_id);
CREATE INDEX IF NOT EXISTS idx_track_points_timestamp ON track_points(timestamp);
`

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, CreateTablesSQL)
	return err
}

// -- Таблица GPS точек
// CREATE TABLE IF NOT EXISTS track_points (
//     id SERIAL PRIMARY KEY,
//     workout_id BIGINT NOT NULL,
//     timestamp TIMESTAMP NOT NULL,
//     latitude DECIMAL(10,8) NOT NULL,
//     longitude DECIMAL(11,8) NOT NULL,
//     elevation DECIMAL(10,2),
//     heart_rate INTEGER,
//     speed DECIMAL(10,2),
//     power INTEGER,
//     cadence INTEGER,
//     temperature DECIMAL(5,2),
//     distance DECIMAL(10,2),
//     grade DECIMAL(5,2),
//     FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
// );

// -- Таблица лапов
// CREATE TABLE IF NOT EXISTS workout_laps (
//     id SERIAL PRIMARY KEY,
//     workout_id BIGINT NOT NULL,
//     lap_number INTEGER NOT NULL,
//     start_time TIMESTAMP NOT NULL,
//     end_time TIMESTAMP NOT NULL,
//     duration INTEGER NOT NULL,
//     distance DECIMAL(10,2),
//     avg_heart_rate INTEGER,
//     max_heart_rate INTEGER,
//     avg_speed DECIMAL(10,2),
//     avg_power INTEGER,
//     avg_cadence INTEGER,
//     calories INTEGER,
//     FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
// );
