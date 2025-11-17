package rest

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Router struct {
	App *fiber.App
	DB  *gorm.DB
}
