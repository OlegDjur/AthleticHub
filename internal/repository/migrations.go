package repository

import "database/sql"

const CreateTablesSQL = `
-- Таблица тренировок
CREATE TABLE IF NOT EXISTS workouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    sub_type TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    duration INTEGER NOT NULL,
    moving_time INTEGER,
    distance REAL,
    avg_speed REAL,
    max_speed REAL,
    avg_pace REAL,
    avg_heart_rate INTEGER,
    max_heart_rate INTEGER,
    resting_hr INTEGER,
    elevation_gain REAL,
    elevation_loss REAL,
    min_elevation REAL,
    max_elevation REAL,
    avg_power INTEGER,
    max_power INTEGER,
    normalized_power INTEGER,
    avg_cadence INTEGER,
    max_cadence INTEGER,
    calories INTEGER,
    avg_temperature REAL,
    min_temperature REAL,
    max_temperature REAL,
    description TEXT,
    device_name TEXT,
    file_type TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Таблица GPS точек
CREATE TABLE IF NOT EXISTS track_points (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_id INTEGER NOT NULL,
    timestamp DATETIME NOT NULL,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    elevation REAL,
    heart_rate INTEGER,
    speed REAL,
    power INTEGER,
    cadence INTEGER,
    temperature REAL,
    distance REAL,
    grade REAL,
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
);

-- Таблица лапов
CREATE TABLE IF NOT EXISTS workout_laps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_id INTEGER NOT NULL,
    lap_number INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    duration INTEGER NOT NULL,
    distance REAL,
    avg_heart_rate INTEGER,
    max_heart_rate INTEGER,
    avg_speed REAL,
    avg_power INTEGER,
    avg_cadence INTEGER,
    calories INTEGER,
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE
);

-- Индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_workouts_start_time ON workouts(start_time);
CREATE INDEX IF NOT EXISTS idx_workouts_type ON workouts(type);
CREATE INDEX IF NOT EXISTS idx_track_points_workout_id ON track_points(workout_id);
CREATE INDEX IF NOT EXISTS idx_track_points_timestamp ON track_points(timestamp);
`

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(CreateTablesSQL)
	return err
}
