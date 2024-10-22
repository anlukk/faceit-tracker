package config

import (

	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"fmt"
	"log"
)

type Config struct {
	TelegramBotAdmins []string  `required:"false"   envconfig:"TELEGRAM_BOT_ADMINS"`
	TelegramAPIURL    string   	`required:"false"   envconfig:"TELEGRAM_API_URL"`
	TelegramToken     string		`required:"true"   envconfig:"TELEGRAM_TOKEN"`
	DatabaseURL       string   	`required:"false"   envconfig:"DATABASE_URL"`
	FaceitAPIToken    string 		`required:"false"   envconfig:"FACEIT_API_TOKEN"`
	FaceitAPISecret   string		`required:"false"   envconfig:"FACEIT_API_SECRET"`
	FaceitAPIURL      string		`required:"false"   envconfig:"FACEIT_API_URL"`
}

func NewConfig() (*Config, error) {
	var (
		newConfig Config
	)

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

  envPath := filepath.Join(filepath.Dir(wd), ".env")

	_, err = os.Stat(envPath)
	if err != nil {
			if os.IsNotExist(err) {
					log.Printf("Warning: .env file not found at %s", envPath)
			} else {
					return nil, err
			}
	}

	err = godotenv.Load(envPath)
	if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	configErr := envconfig.Process("", &newConfig)
	if configErr != nil {
		return nil, fmt.Errorf("config error: %w", configErr)
	}

	return &newConfig, nil
}
