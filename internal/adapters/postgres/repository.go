package postgres

import (
	"context"
	"fmt"
	"workout/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

func NewPostgresAdapter(cfg config.DatabaseConfig) (*postgres, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	// Создаем конфигурацию пула соединений
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга строки подключения: %w", err)
	}

	// Настраиваем пул соединений
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	// Создаем пул соединений
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула соединений: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	return &postgres{db: pool}, nil
}

// Close закрывает соединение с базой данных
func (p *postgres) Close() error {
	if p.db != nil {
		p.db.Close()
	}
	return nil
}
