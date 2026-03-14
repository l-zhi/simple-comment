package config

import (
	"os"
)

type Config struct {
	HTTPPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:   getenv("APP_PORT", "8081"),
		DBHost:     getenv("DB_HOST", "mysql"),
		DBPort:     getenv("DB_PORT", "3306"),
		DBUser:     getenv("DB_USER", "root"),
		DBPassword: getenv("DB_PASSWORD", "password"),
		DBName:     getenv("DB_NAME", "simple_comment"),
	}
	return cfg, nil
}

