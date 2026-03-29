package dto

import (
	"time"

	"github.com/google/uuid"
)

// ─────────────────────────────────────────────
// Requests
// ─────────────────────────────────────────────

type CreateProductRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  uuid.UUID `json:"categoryId"`
	ImageUrl    string    `json:"imageUrl"`
}

type UpdateProductRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  uuid.UUID `json:"categoryId"`
	ImageUrl    string    `json:"imageUrl"`
}

// GetProductsRequest holds optional filters for listing products.
// Pointer fields distinguish "not provided" from zero value —
// nil means "no filter", matching the ProductFilter in the repository.
type GetProductsRequest struct {
	SellerID   *uuid.UUID `json:"sellerId"`
	CategoryID *uuid.UUID `json:"categoryId"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ─────────────────────────────────────────────
// Responses
// ─────────────────────────────────────────────

type CategoryResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type ProductResponse struct {
	Uuid        uuid.UUID        `json:"uuid"`
	SellerID    uuid.UUID        `json:"sellerId"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
	Category    CategoryResponse `json:"category"`
	ImageUrl    string           `json:"imageUrl"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

type GetProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}
