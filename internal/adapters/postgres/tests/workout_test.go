package tests

import (
	"context"
	"testing"
	"time"

	"workout/internal/adapters"
	"workout/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkoutRepository_CreateWorkout(t *testing.T) {
	// Подключение к тестовой БД
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://athletic_user:athletic_password_2024@localhost:5432/athletic_hub")
	require.NoError(t, err)
	defer pool.Close()

	// Создание репозитория
	repo := adapters.NewWorkoutRepository(pool)

	// Тестовые данные
	workout := &entity.Workout{
		UserID:       1, // тестовый пользователь
		Name:         "Тестовая тренировка",
		SportType:    "running",
		Date:         time.Now().Format("2006-01-02"),
		Duration:     "3600",  // 1 час в секундах
		Distance:     "10.50", // 10.5 км
		AvgPace:      "5.5",
		AvgHeartRate: 150,
		MaxHeartRate: 180,
		AvgCadence:   170,
		Calories:     800,
		Description:  "Тестовая тренировка",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Выполнение теста
	err = repo.CreateWorkout(ctx, workout)

	// Проверки
	assert.NoError(t, err)
	assert.NotEmpty(t, workout.ID, "ID должен быть установлен после создания")
}

func TestWorkoutRepository_CreateWorkout_EmptyData(t *testing.T) {
	// Подключение к тестовой БД
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://athletic_user:athletic_password_2024@localhost:5432/athletic_hub")
	require.NoError(t, err)
	defer pool.Close()

	// Создание репозитория
	repo := adapters.NewWorkoutRepository(pool)

	// Тестовые данные с минимальными полями
	workout := &entity.Workout{
		UserID:       1, // тестовый пользователь
		Name:         "Минимальная тренировка",
		SportType:    "walking",
		Date:         time.Now().Format("2006-01-02"),
		Duration:     "1800", // 30 минут
		Distance:     "2.00", // 2 км
		AvgPace:      "6.0",  // темп в мин/км
		AvgHeartRate: 120,
		MaxHeartRate: 140,
		AvgCadence:   150,
		Calories:     200,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Выполнение теста
	err = repo.CreateWorkout(ctx, workout)

	// Проверки
	assert.NoError(t, err)
	assert.NotZero(t, workout.ID, "ID должен быть установлен после создания")
}
