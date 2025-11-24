package rest

import (
	"ecommerce/config"
	"ecommerce/internal/helper"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config config.AppConfig
}
