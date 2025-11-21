package rest

import (
	"ecommerce/internal/helper"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Router struct {
	App  *fiber.App
	DB   *gorm.DB
	Auth helper.Auth
}
