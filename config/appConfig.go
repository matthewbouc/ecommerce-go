package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort     string
	GrpcPort       string
	DatabaseConfig string
	AuthSecret     string
	TwilioToken    string
	TwilioAccount  string
	TwilioNumber   string
}

func SetupEnv() (cfg AppConfig, err error) {

	if os.Getenv("APP_ENV") == "dev" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			return AppConfig{}, err
		}
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return AppConfig{}, errors.New("no HTTP_PORT environment variable set")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = ":50051"
	}

	databaseConfig := os.Getenv("DATABASE_CONFIG")
	if len(databaseConfig) < 1 {
		return AppConfig{}, errors.New("no DATA_SOURCE_NAME environment variable set")
	}

	authSecret := os.Getenv("AUTH_SECRET")
	if len(authSecret) < 1 {
		return AppConfig{}, errors.New("no AUTH_SECRET environment variable set")
	}

	twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if len(twilioToken) < 1 {
		return AppConfig{}, errors.New("no TWILIO_AUTH_TOKEN environment variable set")
	}

	twilioAccount := os.Getenv("TWILIO_ACCOUNT_SID")
	if len(twilioAccount) < 1 {
		return AppConfig{}, errors.New("no TWILIO_ACCOUNT environment variable set")
	}

	myPhoneNumber := os.Getenv("TWILIO_NUMBER")
	if len(myPhoneNumber) < 1 {
		return AppConfig{}, errors.New("no MY_NUMBER environment variable set")
	}

	return AppConfig{
		ServerPort:     httpPort,
		GrpcPort:       grpcPort,
		DatabaseConfig: databaseConfig,
		AuthSecret:     authSecret,
		TwilioAccount:  twilioAccount,
		TwilioToken:    twilioToken,
		TwilioNumber:   myPhoneNumber,
	}, nil
}
