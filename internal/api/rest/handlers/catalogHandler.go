package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/helper"
	"ecommerce/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CatalogHandler struct {
	svc  service.CatalogService
	auth helper.Auth
}

func NewCatalogHandler(svc service.CatalogService, auth helper.Auth) CatalogHandler {
	return CatalogHandler{svc: svc, auth: auth}
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	// Use the shared CatalogService instance from RestHandler — the same
	// instance is also used by the gRPC server, so there's only one.
	h := NewCatalogHandler(rh.CatalogSvc, rh.Auth)

	catalog := app.Group("/catalog")

	// ─── Public — no auth required ────────────────────────────────────────
	catalog.Get("/products", h.GetProducts)
	catalog.Get("/products/:id", h.GetProduct)
	catalog.Get("/categories", h.GetCategories)

	// ─── Seller only ───────────────────────────────────────────────────────
	// Authorize validates JWT; RequireRole(SELLER) rejects non-sellers.
	seller := catalog.Group("/", rh.Auth.Authorize, rh.Auth.RequireRole("seller"))
	seller.Post("/products", h.CreateProduct)
	seller.Put("/products/:id", h.UpdateProduct)
	seller.Delete("/products/:id", h.DeleteProduct)

	// ─── Admin / internal — seller-gated for now ──────────────────────────
	// Categories are managed by admins in production; seller-role is a
	// temporary stand-in until an admin role is introduced.
	seller.Post("/categories", h.CreateCategory)
}

// ─────────────────────────────────────────────
// Product handlers
// ─────────────────────────────────────────────

func (h *CatalogHandler) GetProducts(ctx fiber.Ctx) error {
	// Query params are optional — nil pointer fields mean "no filter".
	// Fiber v3 returns query params as strings; parse with a fallback default.
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	req := dto.GetProductsRequest{
		Page:     page,
		PageSize: pageSize,
	}

	if sellerStr := ctx.Query("sellerId"); sellerStr != "" {
		sellerUuid, err := uuid.Parse(sellerStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid sellerId",
			})
		}
		req.SellerID = &sellerUuid
	}

	if catStr := ctx.Query("categoryId"); catStr != "" {
		catUuid, err := uuid.Parse(catStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid categoryId",
			})
		}
		req.CategoryID = &catUuid
	}

	result, err := h.svc.GetProducts(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch products",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "success",
		"products": result.Products,
		"total":    result.Total,
	})
}

func (h *CatalogHandler) GetProduct(ctx fiber.Ctx) error {
	productUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid product id",
		})
	}

	product, err := h.svc.GetProduct(productUuid)
	if err != nil {
		// Surface a 404 when GORM returns record-not-found.
		if isNotFound(err) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "product not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch product",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"product": product,
	})
}

func (h *CatalogHandler) CreateProduct(ctx fiber.Ctx) error {
	seller := h.auth.GetCurrentUser(ctx)

	var req dto.CreateProductRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	product, err := h.svc.CreateProduct(seller.Uuid, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "product created",
		"product": product,
	})
}

func (h *CatalogHandler) UpdateProduct(ctx fiber.Ctx) error {
	seller := h.auth.GetCurrentUser(ctx)

	productUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid product id",
		})
	}

	var req dto.UpdateProductRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	product, err := h.svc.UpdateProduct(seller.Uuid, productUuid, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "product updated",
		"product": product,
	})
}

func (h *CatalogHandler) DeleteProduct(ctx fiber.Ctx) error {
	seller := h.auth.GetCurrentUser(ctx)

	productUuid, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid product id",
		})
	}

	if err := h.svc.DeleteProduct(seller.Uuid, productUuid); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "product deleted",
	})
}

// ─────────────────────────────────────────────
// Category handlers
// ─────────────────────────────────────────────

func (h *CatalogHandler) GetCategories(ctx fiber.Ctx) error {
	categories, err := h.svc.GetCategories()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to fetch categories",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "success",
		"categories": categories,
	})
}

func (h *CatalogHandler) CreateCategory(ctx fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	category, err := h.svc.CreateCategory(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "category created",
		"category": category,
	})
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

// isNotFound checks if a GORM error is a record-not-found error so handlers
// can return 404 instead of 500.
func isNotFound(err error) bool {
	return err != nil && err.Error() == gorm.ErrRecordNotFound.Error()
}
