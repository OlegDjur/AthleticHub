package config

import "time"

// Config основная структура конфигурации
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	//Auth     AuthConfig
	//Storage  StorageConfig
	//Logging  LoggingConfig
}

// ServerConfig настройки сервера
type ServerConfig struct {
	Port         int
	Host         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Debug        bool
}

// DatabaseConfig настройки базы данных
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}
