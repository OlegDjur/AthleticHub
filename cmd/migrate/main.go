package main

import (
	"context"
	"log"
	"workout/internal/adapters/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// Подключение к БД
	pool, err := pgxpool.New(ctx, "postgres://athletic_user:athletic_password_2024@localhost:5432/athletic_hub")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer pool.Close()

	// Запуск миграций
	log.Println("Запуск миграций...")
	err = postgres.RunMigrations(ctx, pool)
	if err != nil {
		log.Fatal("Ошибка выполнения миграций:", err)
	}

	log.Println("Миграции выполнены успешно!")
}
