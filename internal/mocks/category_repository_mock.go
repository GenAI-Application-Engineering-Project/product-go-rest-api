package mocks

import (
	"context"

	"product-services/internal/models"
	"product-services/internal/shared"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetCategoryByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) ListCategories(
	ctx context.Context,
	opts shared.ListOptions,
) (*models.ListCategoriesResult, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*models.ListCategoriesResult), args.Error(1)
}

func (m *MockCategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
