package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port string
		Host string
	}
	Database struct {
		Path string
	}
	Upload struct {
		MaxFileSize  int64 // байты
		AllowedTypes []string
	}
}

func LoadConfig() *Config {
	cfg := &Config{}

	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.Host = getEnv("HOST", "localhost")
	cfg.Database.Path = getEnv("DB_PATH", "./workout.db")

	maxSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "10485760"), 10, 64) // 10MB
	cfg.Upload.MaxFileSize = maxSize
	cfg.Upload.AllowedTypes = []string{".fit", ".tcx", ".gpx"}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
