package repository

import (
	"context"
	"ecommerce/internal/domain"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductFilter holds optional query parameters for listing products.
// Pointer fields mean nil = "no filter applied" — avoids overloading
// the function signature with many optional args.
type ProductFilter struct {
	SellerID   *uuid.UUID
	CategoryID *uuid.UUID
	Page       int
	PageSize   int
}

type CatalogRepository interface {
	// Product operations
	CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error)
	GetProductByUuid(ctx context.Context, productUuid uuid.UUID) (*domain.Product, error)
	GetProducts(ctx context.Context, filter ProductFilter) ([]domain.Product, int64, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, product *domain.Product) error

	// Category operations
	CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
	GetCategoryByUuid(ctx context.Context, categoryUuid uuid.UUID) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]domain.Category, error)
}

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{db: db}
}

// ─────────────────────────────────────────────
// Product
// ─────────────────────────────────────────────

func (r *catalogRepository) CreateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	if err := r.db.WithContext(ctx).Create(&product).Error; err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}
	return product, nil
}

func (r *catalogRepository) GetProductByUuid(ctx context.Context, productUuid uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	// Preload("Category") fetches the related Category row in a single extra query, so it's not empty struct
	err := r.db.WithContext(ctx).Preload("Category").
		Where("uuid = ?", productUuid).
		First(&product).Error
	if err != nil {
		return nil, fmt.Errorf("get product by uuid %s: %w", productUuid, err)
	}
	return &product, nil
}

func (r *catalogRepository) GetProducts(ctx context.Context, filter ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{})

	if filter.SellerID != nil {
		query = query.Where("seller_id = ?", *filter.SellerID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	// Count the total matching rows before applying LIMIT/OFFSET.
	// This gives the caller the information needed to calculate total pages.
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	// Apply safe defaults so callers don't have to think about edge cases.
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize

	err := query.
		Preload("Category").
		Limit(filter.PageSize).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		return nil, 0, fmt.Errorf("get products: %w", err)
	}

	return products, total, nil
}

func (r *catalogRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	// clause.Returning{} updates the pointer in-place with the DB's version of
	// the record after the update — avoids a separate SELECT to confirm changes.
	err := r.db.WithContext(ctx).Model(product).
		Clauses(clause.Returning{}).
		Where("uuid = ?", product.Uuid).
		Updates(product).Error
	if err != nil {
		return fmt.Errorf("update product %s: %w", product.Uuid, err)
	}
	return nil
}

func (r *catalogRepository) DeleteProduct(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Delete(product).Error; err != nil {
		return fmt.Errorf("delete product %s: %w", product.Uuid, err)
	}
	return nil
}

// ─────────────────────────────────────────────
// Category
// ─────────────────────────────────────────────

func (r *catalogRepository) CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return domain.Category{}, fmt.Errorf("create category: %w", err)
	}
	return category, nil
}

func (r *catalogRepository) GetCategoryByUuid(ctx context.Context, categoryUuid uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.WithContext(ctx).Where("uuid = ?", categoryUuid).First(&category).Error; err != nil {
		return nil, fmt.Errorf("get category by uuid %s: %w", categoryUuid, err)
	}
	return &category, nil
}

func (r *catalogRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	return categories, nil
}
