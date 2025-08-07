package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

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

// LoadConfig загружает конфигурацию из файла .env
func LoadConfig(envPath string) (*Config, error) {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	config := &Config{}

	// Загружаем настройки сервера
	config.Server = ServerConfig{
		Port:         getEnvAsInt("SERVER_PORT", 8080),
		Host:         getEnv("SERVER_HOST", "localhost"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		Debug:        getEnvAsBool("DEBUG", false),
	}

	// Загружаем настройки базы данных
	config.Database = DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		Name:            getEnv("DB_NAME", "workout"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}

	return config, nil
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает переменную окружения как int или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool получает переменную окружения как bool или возвращает значение по умолчанию
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsDuration получает переменную окружения как time.Duration или возвращает значение по умолчанию
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
