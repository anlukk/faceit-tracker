package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotAdmins []string `env:"TELEGRAM_BOT_ADMINS"`
	TelegramAPIURL    string   `env:"TELEGRAM_API_URL"`
	TelegramToken     string   `env:"TELEGRAM_TOKEN"`
	DatabaseURL       string   `env:"DATABASE_URL"`

	// Faceit API
	FaceitAPIToken string `env:"FACEIT_API_TOKEN"`
	FaceitAPIURL   string `env:"FACEIT_API_URL"`

	// Logger
	LoggerLevel string `env:"LOGGER_LEVEL"`

	// Database
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBSSLMode  string `env:"SSL_MODE"`

	MaxIdleConns    int `env:"MAX_IDLE_CONNS"`
	MaxOpenConns    int `env:"MAX_OPEN_CONNS"`
	ConnMaxLifetime int `env:"CONN_MAX_LIFETIME"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return &cfg, nil
}
