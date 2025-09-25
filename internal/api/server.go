package api

import (
	"ecommerce/config"
	"log"

	"github.com/gofiber/fiber/v3"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	app.Get("/health", Healthcheck)

	log.Fatal(app.Listen(config.ServerPort))
}

func Healthcheck(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"healthy": true,
	})
}
