package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {

	app := rh.App

	userSvc := service.UserService{
		UserRepository:        repository.NewUserRepository(rh.DB),
		BankAccountRepository: repository.NewBankAccountRepository(rh.DB),
		CartRepository:        repository.NewCartRepository(rh.DB),
		OrderRepository:       repository.NewOrderRepository(rh.DB),
		Auth:                  rh.Auth,
		Config:                rh.Config,
		SmsClient:             rh.SmsClient,
	}
	userHandler := UserHandler{
		service: userSvc,
	}

	user := app.Group("/user")

	// #### public ####
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)

	// #### private ####
	// TODO, eventually want to verify role and that user is not deleted. Add additional handlers for that logic?
	pvtUser := user.Group("/", rh.Auth.Authorize)

	pvtUser.Delete("/", userHandler.DeleteUser)

	pvtUser.Get("/verify", userHandler.GetVerificationCode)
	pvtUser.Post("/verify", userHandler.Verify)

	pvtUser.Get("/profile", userHandler.GetProfile)
	pvtUser.Post("/profile", userHandler.UpdateProfile)

	pvtUser.Get("/cart", userHandler.GetCart)
	pvtUser.Post("/cart", userHandler.AddToCart)

	pvtUser.Get("/order", userHandler.GetOrders)
	pvtUser.Get("/order/:id", userHandler.GetOrder)

	pvtUser.Post("/become-seller", userHandler.BecomeSeller)
}

func (h *UserHandler) Register(ctx fiber.Ctx) error {

	user := dto.RegisterRequest{}

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
	loginAttempt := dto.LoginRequest{}

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

	user := h.service.Auth.GetCurrentUser(ctx)

	code, err := h.service.GetVerificationCode(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "get verification code successful",
		"code":    code,
	})
}

func (h *UserHandler) Verify(ctx fiber.Ctx) error {

	user := h.service.Auth.GetCurrentUser(ctx)

	var req dto.VerificationCodeInput

	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	err := h.service.VerifyCode(user.Uuid, req.Code)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "code successfully verified",
	})
}

func (h *UserHandler) GetProfile(ctx fiber.Ctx) error {
	user := h.service.Auth.GetCurrentUser(ctx)
	profile, err := h.service.GetProfile(user.Uuid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"profile": profile,
	})
}

func (h *UserHandler) UpdateProfile(ctx fiber.Ctx) error {
	user := h.service.Auth.GetCurrentUser(ctx)

	var req dto.UpdateProfileRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	if err := h.service.UpdateProfile(user.Uuid, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "profile updated",
	})
}

func (h *UserHandler) GetCart(ctx fiber.Ctx) error {
	user := h.service.Auth.GetCurrentUser(ctx)
	items, err := h.service.GetCart(user.Uuid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"cart":    items,
	})
}

func (h *UserHandler) AddToCart(ctx fiber.Ctx) error {
	user := h.service.Auth.GetCurrentUser(ctx)

	var req dto.AddToCartRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	if err := h.service.AddToCart(user.Uuid, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "item added to cart",
	})
}

func (h *UserHandler) GetOrders(ctx fiber.Ctx) error {
	user := h.service.Auth.GetCurrentUser(ctx)
	orders, err := h.service.GetOrders(user.Uuid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"orders":  orders,
	})
}

func (h *UserHandler) GetOrder(ctx fiber.Ctx) error {
	user := h.service.Auth.GetCurrentUser(ctx)

	orderUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order id",
		})
	}

	order, err := h.service.GetOrderById(user.Uuid, orderUuid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"order":   order,
	})
}

func (h *UserHandler) BecomeSeller(ctx fiber.Ctx) error {

	user := h.service.Auth.GetCurrentUser(ctx)

	req := dto.BecomeSellerRequest{}
	err := ctx.Bind().Body(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	req.Uuid = user.Uuid
	token, err := h.service.BecomeSeller(req)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to register as a seller",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully registered as a seller",
		"token":   token,
	})
}
