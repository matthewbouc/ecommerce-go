package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// UserHandler handles HTTP requests for the /user routes.
// svc holds the business logic interface; auth handles context extraction.
type UserHandler struct {
	svc  service.UserService
	auth helper.Auth
}

func NewUserHandler(svc service.UserService, auth helper.Auth) UserHandler {
	return UserHandler{svc: svc, auth: auth}
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App

	userSvc := service.NewUserService(
		repository.NewUserRepository(rh.DB),
		repository.NewBankAccountRepository(rh.DB),
		repository.NewCartRepository(rh.DB),
		repository.NewOrderRepository(rh.DB),
		rh.Auth,
		rh.Config,
		rh.SmsClient,
		rh.CatalogClient,
	)
	userHandler := NewUserHandler(userSvc, rh.Auth)

	user := app.Group("/user")

	// #### public ####
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)

	// #### private ####
	// Authorize validates the JWT; requireActiveUser ensures the account hasn't been soft-deleted.
	pvtUser := user.Group("/", rh.Auth.Authorize, userHandler.requireActiveUser)

	pvtUser.Delete("/", userHandler.DeleteUser)

	pvtUser.Get("/verify", userHandler.GetVerificationCode)
	pvtUser.Post("/verify", userHandler.Verify)

	pvtUser.Get("/profile", userHandler.GetProfile)
	pvtUser.Post("/profile", userHandler.UpdateProfile)

	pvtUser.Get("/cart", userHandler.GetCart)
	pvtUser.Post("/cart", userHandler.AddToCart)
	pvtUser.Delete("/cart/:id", userHandler.RemoveFromCart)

	pvtUser.Get("/order", userHandler.GetOrders)
	pvtUser.Get("/order/:id", userHandler.GetOrder)
	pvtUser.Post("/checkout", userHandler.Checkout)
	pvtUser.Patch("/order/:id/status", userHandler.UpdateOrderStatus)

	// become-seller is restricted to buyers; sellers cannot call it twice.
	pvtUser.Post("/become-seller", rh.Auth.RequireRole(domain.BUYER), userHandler.BecomeSeller)
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

	token, err := h.svc.Register(user)
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

	token, err := h.svc.Login(loginAttempt)
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
	user := h.auth.GetCurrentUser(ctx)
	err := h.svc.DeleteUser(user.Uuid)
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
	user := h.auth.GetCurrentUser(ctx)

	code, err := h.svc.GetVerificationCode(user)
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
	user := h.auth.GetCurrentUser(ctx)

	var req dto.VerificationCodeInput
	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	if err := h.svc.VerifyCode(user.Uuid, req.Code); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "code successfully verified",
	})
}

func (h *UserHandler) GetProfile(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)
	profile, err := h.svc.GetProfile(user.Uuid)
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
	user := h.auth.GetCurrentUser(ctx)

	var req dto.UpdateProfileRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	if err := h.svc.UpdateProfile(user.Uuid, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "profile updated",
	})
}

func (h *UserHandler) GetCart(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)
	items, err := h.svc.GetCart(user.Uuid)
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
	user := h.auth.GetCurrentUser(ctx)

	var req dto.AddToCartRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	if err := h.svc.AddToCart(user.Uuid, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "item added to cart",
	})
}

func (h *UserHandler) GetOrders(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)
	orders, err := h.svc.GetOrders(user.Uuid)
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
	user := h.auth.GetCurrentUser(ctx)

	orderUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order id",
		})
	}

	order, err := h.svc.GetOrderById(user.Uuid, orderUuid)
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

func (h *UserHandler) Checkout(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)

	order, err := h.svc.Checkout(user.Uuid)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "order created",
		"order":   order,
	})
}

func (h *UserHandler) UpdateOrderStatus(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)

	orderUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order id",
		})
	}

	var req dto.UpdateOrderStatusRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	if err := h.svc.UpdateOrderStatus(user.Uuid, orderUuid, req.Status); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "order status updated",
	})
}

// requireActiveUser rejects requests from soft-deleted accounts.
// It runs after Authorize so ctx.Locals("user") is already populated.
func (h *UserHandler) requireActiveUser(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)
	if !h.svc.IsActiveUser(user.Uuid) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "account not found or has been deleted",
		})
	}
	return ctx.Next()
}

func (h *UserHandler) RemoveFromCart(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)

	itemId, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid cart item id",
		})
	}

	if err := h.svc.RemoveFromCart(user.Uuid, uint(itemId)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "item removed from cart",
	})
}

func (h *UserHandler) BecomeSeller(ctx fiber.Ctx) error {
	user := h.auth.GetCurrentUser(ctx)

	req := dto.BecomeSellerRequest{}
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	req.Uuid = user.Uuid
	token, err := h.svc.BecomeSeller(req)
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
