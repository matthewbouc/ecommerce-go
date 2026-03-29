package service

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type CatalogService struct {
	CatalogRepository repository.CatalogRepository
}

// ─────────────────────────────────────────────
// Products
// ─────────────────────────────────────────────

func (s CatalogService) CreateProduct(sellerUuid uuid.UUID, input dto.CreateProductRequest) (dto.ProductResponse, error) {
	product := domain.Product{
		SellerID:    sellerUuid,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		CategoryID:  input.CategoryID,
		ImageUrl:    input.ImageUrl,
	}

	if err := product.Validate(); err != nil {
		return dto.ProductResponse{}, err
	}

	created, err := s.CatalogRepository.CreateProduct(product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Fetch the full record so the response includes the populated Category association.
	// CreateProduct returns the saved product but GORM doesn't preload associations on insert.
	full, err := s.CatalogRepository.GetProductByUuid(created.Uuid)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return toProductResponse(*full), nil
}

func (s CatalogService) GetProduct(productUuid uuid.UUID) (dto.ProductResponse, error) {
	product, err := s.CatalogRepository.GetProductByUuid(productUuid)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return toProductResponse(*product), nil
}

func (s CatalogService) GetProducts(req dto.GetProductsRequest) (dto.GetProductsResponse, error) {
	products, total, err := s.CatalogRepository.GetProducts(repository.ProductFilter{
		SellerID:   req.SellerID,
		CategoryID: req.CategoryID,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		return dto.GetProductsResponse{}, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		response[i] = toProductResponse(p)
	}

	return dto.GetProductsResponse{
		Products: response,
		Total:    total,
	}, nil
}

func (s CatalogService) UpdateProduct(sellerUuid uuid.UUID, productUuid uuid.UUID, input dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product, err := s.CatalogRepository.GetProductByUuid(productUuid)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Authorization: only the seller who owns the product can update it.
	// This check belongs here in the service — it's a business rule, not an HTTP concern.
	if product.SellerID != sellerUuid {
		return dto.ProductResponse{}, errors.New("unauthorized: product does not belong to this seller")
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.CategoryID = input.CategoryID
	product.ImageUrl = input.ImageUrl

	if err := product.Validate(); err != nil {
		return dto.ProductResponse{}, err
	}

	if err := s.CatalogRepository.UpdateProduct(product); err != nil {
		return dto.ProductResponse{}, err
	}

	// UpdateProduct uses clause.Returning{} so product is updated in-place by GORM.
	return toProductResponse(*product), nil
}

func (s CatalogService) DeleteProduct(sellerUuid uuid.UUID, productUuid uuid.UUID) error {
	product, err := s.CatalogRepository.GetProductByUuid(productUuid)
	if err != nil {
		return err
	}

	// Authorization: only the owning seller can delete their product.
	if product.SellerID != sellerUuid {
		return errors.New("unauthorized: product does not belong to this seller")
	}

	return s.CatalogRepository.DeleteProduct(product)
}

// ─────────────────────────────────────────────
// Categories
// ─────────────────────────────────────────────

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	category := domain.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := category.Validate(); err != nil {
		return dto.CategoryResponse{}, err
	}

	created, err := s.CatalogRepository.CreateCategory(category)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	return toCategoryResponse(created), nil
}

func (s CatalogService) GetCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.CatalogRepository.GetCategories()
	if err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i, c := range categories {
		response[i] = toCategoryResponse(c)
	}
	return response, nil
}

// ─────────────────────────────────────────────
// Helpers — domain → DTO mapping
// ─────────────────────────────────────────────

func toProductResponse(p domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		Uuid:        p.Uuid,
		SellerID:    p.SellerID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    toCategoryResponse(p.Category),
		ImageUrl:    p.ImageUrl,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toCategoryResponse(c domain.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		Uuid:        c.Uuid,
		Name:        c.Name,
		Description: c.Description,
	}
}
