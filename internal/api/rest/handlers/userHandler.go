package handlers

import (
	"ecommerce/internal/api/rest"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	// UserService
}

func SetupUserRoutes(router *rest.Router) {

	app := router.App
	handler := UserHandler{}

	// public endpoints
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	// private endpoints
	app.Get("/verify", handler.GetVerificationCode)
	app.Post("/verify", handler.Verify)
	app.Get("/profile", handler.GetProfile)
	app.Post("/profile", handler.CreateProfile)

	app.Get("/cart", handler.GetCart)
	app.Post("/cart", handler.AddToCart)
	app.Get("/order", handler.GetOrders)
	app.Get("/order/:id", handler.GetOrder)

	app.Post("/become-seller", handler.BecomeSeller)
}

func (h *UserHandler) Register(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "register success",
	})
}
func (h *UserHandler) Login(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login success",
	})
}
func (h *UserHandler) GetVerificationCode(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) Verify(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) GetProfile(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) CreateProfile(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) AddToCart(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) GetCart(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) GetOrders(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) GetOrder(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
func (h *UserHandler) BecomeSeller(ctx fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
