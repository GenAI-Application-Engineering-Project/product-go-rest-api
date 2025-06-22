package mocks

import (
	"context"
	"product-services/internal/models"
	"product-services/internal/shared"

	"github.com/google/uuid"
)

type MockCategoryRepository struct {
	GetCategoryByIDFunc    func(ctx context.Context, id uuid.UUID) (*models.Category, error)
	ListCategoriesFunc     func(ctx context.Context, opts shared.ListOptions) (*models.ListCategoriesResult, error)
	CreateCategoryFunc     func(ctx context.Context, category *models.Category) error
	UpdateCategoryFunc     func(ctx context.Context, category *models.Category) error
	DeleteCategoryFunc     func(ctx context.Context, id uuid.UUID) error
}

func (m *MockCategoryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	return m.GetCategoryByIDFunc(ctx, id)
}

func (m *MockCategoryRepository) ListCategories(ctx context.Context, opts shared.ListOptions) (*models.ListCategoriesResult, error) {
	return m.ListCategoriesFunc(ctx, opts)
}

func (m *MockCategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	return m.CreateCategoryFunc(ctx, category)
}

func (m *MockCategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	return m.UpdateCategoryFunc(ctx, category)
}

func (m *MockCategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return m.DeleteCategoryFunc(ctx, id)
}