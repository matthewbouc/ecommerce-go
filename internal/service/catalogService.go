package service

import (
	"ecommerce/internal/domain"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"errors"

	"github.com/google/uuid"
)

// CatalogService is the public interface consumed by the REST handler and
// the gRPC server. Both call the same business logic through this contract.
type CatalogService interface {
	CreateProduct(sellerUuid uuid.UUID, input dto.CreateProductRequest) (dto.ProductResponse, error)
	GetProduct(productUuid uuid.UUID) (dto.ProductResponse, error)
	GetProducts(req dto.GetProductsRequest) (dto.GetProductsResponse, error)
	UpdateProduct(sellerUuid uuid.UUID, productUuid uuid.UUID, input dto.UpdateProductRequest) (dto.ProductResponse, error)
	DeleteProduct(sellerUuid uuid.UUID, productUuid uuid.UUID) error
	CreateCategory(input dto.CreateCategoryRequest) (dto.CategoryResponse, error)
	GetCategories() ([]dto.CategoryResponse, error)
}

// catalogServiceImpl is the concrete implementation. Unexported — callers
// receive a CatalogService interface from NewCatalogService.
type catalogServiceImpl struct {
	repo repository.CatalogRepository
}

func NewCatalogService(repo repository.CatalogRepository) CatalogService {
	return &catalogServiceImpl{repo: repo}
}

// ─────────────────────────────────────────────
// Products
// ─────────────────────────────────────────────

func (s *catalogServiceImpl) CreateProduct(sellerUuid uuid.UUID, input dto.CreateProductRequest) (dto.ProductResponse, error) {
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

	created, err := s.repo.CreateProduct(product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Fetch the full record so the response includes the populated Category association.
	// CreateProduct returns the saved product but GORM doesn't preload associations on insert.
	full, err := s.repo.GetProductByUuid(created.Uuid)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return toProductResponse(*full), nil
}

func (s *catalogServiceImpl) GetProduct(productUuid uuid.UUID) (dto.ProductResponse, error) {
	product, err := s.repo.GetProductByUuid(productUuid)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return toProductResponse(*product), nil
}

func (s *catalogServiceImpl) GetProducts(req dto.GetProductsRequest) (dto.GetProductsResponse, error) {
	products, total, err := s.repo.GetProducts(repository.ProductFilter{
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

func (s *catalogServiceImpl) UpdateProduct(sellerUuid uuid.UUID, productUuid uuid.UUID, input dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product, err := s.repo.GetProductByUuid(productUuid)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Authorization: only the seller who owns the product can update it.
	// This check belongs in the service — it's a business rule, not an HTTP concern.
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

	if err := s.repo.UpdateProduct(product); err != nil {
		return dto.ProductResponse{}, err
	}

	// UpdateProduct uses clause.Returning{} so product is updated in-place by GORM.
	return toProductResponse(*product), nil
}

func (s *catalogServiceImpl) DeleteProduct(sellerUuid uuid.UUID, productUuid uuid.UUID) error {
	product, err := s.repo.GetProductByUuid(productUuid)
	if err != nil {
		return err
	}

	// Authorization: only the owning seller can delete their product.
	if product.SellerID != sellerUuid {
		return errors.New("unauthorized: product does not belong to this seller")
	}

	return s.repo.DeleteProduct(product)
}

// ─────────────────────────────────────────────
// Categories
// ─────────────────────────────────────────────

func (s *catalogServiceImpl) CreateCategory(input dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	category := domain.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := category.Validate(); err != nil {
		return dto.CategoryResponse{}, err
	}

	created, err := s.repo.CreateCategory(category)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	return toCategoryResponse(created), nil
}

func (s *catalogServiceImpl) GetCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.repo.GetCategories()
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
