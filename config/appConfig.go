package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort     string
	DatabaseConfig string
}

func SetupEnv() (cfg AppConfig, err error) {

	if os.Getenv("APP_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			return AppConfig{}, err
		}
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return AppConfig{}, errors.New("no HTTP_PORT environment variable set")
	}

	databaseConfig := os.Getenv("DATABASE_CONFIG")
	if len(databaseConfig) < 1 {
		return AppConfig{}, errors.New("no DATA_SOURCE_NAME environment variable set")
	}

	return AppConfig{ServerPort: httpPort, DatabaseConfig: databaseConfig}, nil
}
