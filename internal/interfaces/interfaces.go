package interfaces

import (
	"context"
	"time"

	"product-services/internal/models"
	"product-services/internal/shared"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type (
	// CategoryRepository defines methods for CRUD operations on categories.
	CategoryRepository interface {
		GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
		ListCategories(
			ctx context.Context,
			listOptions shared.ListOptions,
		) (*models.ListCategoriesResult, error)
		CreateCategory(ctx context.Context, category *models.Category) error
		UpdateCategory(ctx context.Context, category *models.Category) error
		DeleteCategory(ctx context.Context, id uuid.UUID) error
	}

	// ProductRepository defines methods for CRUD operations on products.
	ProductRepository interface {
		GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
		ListProducts(
			ctx context.Context,
			listOptions shared.ListOptions,
		) (*models.ListProductsResult, error)
		CreateProduct(ctx context.Context, product *models.Product) error
		UpdateProduct(ctx context.Context, product *models.Product) error
		DeleteProduct(ctx context.Context, id uuid.UUID) error
	}

	// AppLogger defines methods for logging
	AppLogger interface {
		// Logger returns the underlying zerolog.Logger instance.
		Logger() zerolog.Logger
		Fatal(err error, msg string)
	}

	// SystemUtil provides access to system-level utilities such as time and UUID generation.
	// Useful for testability and mocking.
	SystemUtil interface {
		CurrentTime() time.Time
		NewUUID() uuid.UUID
	}
)
