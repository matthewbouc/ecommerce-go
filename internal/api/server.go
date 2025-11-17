package api

import (
	"ecommerce/config"
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/domain"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	database, err := gorm.Open(postgres.Open(config.DatabaseConfig), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection fatal error: %v\n", err)
	}
	log.Println("database connection established")

	err = database.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("database migration fatal error: %v\n", err)
	}

	router := &rest.Router{
		App: app,
		DB:  database,
	}

	setupRoutes(router)

	log.Fatal(app.Listen(config.ServerPort))
}

func setupRoutes(router *rest.Router) {
	// user handler
	handlers.SetupUserRoutes(router)
	//catalog handler
	//transactionhandler
}
