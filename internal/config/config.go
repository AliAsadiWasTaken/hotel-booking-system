package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Logger   LoggerConfig
	HTTP     HTTPConfig
}

type LoggerConfig struct {
	Level     string
	Format    string
	AddSource bool
}

type HTTPConfig struct {
	Address string
}

type DatabaseConfig struct {
	DB       string
	Host     string
	Port     string
	User     string
	Password string
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, falling back to system environment variables")
	}

	cfg := Config{
		Database: DatabaseConfig{
			DB:       os.Getenv("POSTGRES_DB"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
		Logger: LoggerConfig{
			Level:     os.Getenv("LOG_LEVEL"),
			Format:    os.Getenv("LOG_FORMAT"),
			AddSource: os.Getenv("LOG_ADD_SOURCE") == "true",
		},
		HTTP: HTTPConfig{
			Address: os.Getenv("HTTP_ADDRESS"),
		},
	}

	if cfg.Database.DB == "" {
		return Config{}, fmt.Errorf("DATABASE_DB environment variable not set")
	}

	if cfg.Database.Host == "" {
		return Config{}, fmt.Errorf("DATABASE_HOST environment variable not set")
	}

	if cfg.Database.Port == "" {
		return Config{}, fmt.Errorf("DATABASE_PORT environment variable not set")
	}

	if cfg.Database.User == "" {
		return Config{}, fmt.Errorf("DATABASE_USER environment variable not set")
	}

	if cfg.Database.Password == "" {
		return Config{}, fmt.Errorf("DATABASE_Password environment variable not set")
	}

	if cfg.Logger.Level == "" {
		return Config{}, fmt.Errorf("LOG_LEVEL environment variable not set")
	}

	if cfg.Logger.Format == "" {
		return Config{}, fmt.Errorf("LOG_FORMAT environment variable not set")
	}

	return cfg, nil
}
