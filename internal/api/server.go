package api

import (
	"ecommerce/config"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/domain"
	"ecommerce/internal/helper"
	"ecommerce/pkg/notification/sms/provider"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {

	database, err := gorm.Open(postgres.Open(config.DatabaseConfig), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection fatal error: %v\n", err)
	}
	log.Println("database connection established")

	err = database.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("database migration fatal error: %v\n", err)
	}

	app := fiber.New()

	auth := helper.SetupAuth(config.AuthSecret)

	smsClient := provider.NewTwilioSmsClient(config)

	restHandler := &rest.RestHandler{
		App:       app,
		DB:        database,
		Auth:      auth,
		Config:    config,
		SmsClient: smsClient,
	}

	setupRoutes(restHandler)

	log.Fatal(
		app.Listen(config.ServerPort),
	)
}

func setupRoutes(rh *rest.RestHandler) {
	// user handler
	handlers.SetupUserRoutes(rh)
	//catalog handler
	//transactionhandler
}
