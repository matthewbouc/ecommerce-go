package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	// UserService
	service service.UserService
}

func SetupUserRoutes(router *rest.Router) {

	app := router.App

	serv := service.UserService{
		UserRepository: repository.NewUserRepository(router.DB),
	}
	handler := UserHandler{
		service: serv,
	}

	// #### public endpoints ####
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	// #### private endpoints ####
	app.Delete("/user/:uuid", handler.DeleteUser)

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

	user := dto.RegisterDTO{}

	// TODO: add some validation here?
	err := ctx.Bind().Body(&user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	token, err := h.service.Register(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error on sign up",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "register success",
		"token":   token,
	})
}
func (h *UserHandler) Login(ctx fiber.Ctx) error {
	loginAttempt := dto.LoginDTO{}

	err := ctx.Bind().Body(&loginAttempt)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide valid login inputs",
		})
	}

	token, err := h.service.Login(loginAttempt)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": token,
	})
}

func (h *UserHandler) DeleteUser(ctx fiber.Ctx) error {
	uuid := ctx.Params("uuid")
	err := h.service.DeleteUser(uuid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted",
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
