package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	service service.UserService
}

func SetupUserRoutes(router *rest.Router) {

	app := router.App

	userSvc := service.UserService{
		UserRepository: repository.NewUserRepository(router.DB),
		Auth:           router.Auth,
	}
	userHandler := UserHandler{
		service: userSvc,
	}

	user := app.Group("/user")

	// #### public ####
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)

	// #### private ####
	pvtUser := user.Group("/", router.Auth.Authorize)

	pvtUser.Delete("/", userHandler.DeleteUser)

	pvtUser.Get("/verify", userHandler.GetVerificationCode)
	pvtUser.Post("/verify", userHandler.Verify)

	pvtUser.Get("/profile", userHandler.GetProfile)
	pvtUser.Post("/profile", userHandler.CreateProfile)

	pvtUser.Get("/cart", userHandler.GetCart)
	pvtUser.Post("/cart", userHandler.AddToCart)

	pvtUser.Get("/order", userHandler.GetOrders)
	pvtUser.Get("/order/:id", userHandler.GetOrder)

	pvtUser.Post("/become-seller", userHandler.BecomeSeller)
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
		"message": "register successful",
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
		"message": "login successful",
		"token":   token,
	})
}

func (h *UserHandler) DeleteUser(ctx fiber.Ctx) error {
	user := h.service.Auth.GetCurrentUser(ctx)
	err := h.service.DeleteUser(user.Uuid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user deleted",
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
