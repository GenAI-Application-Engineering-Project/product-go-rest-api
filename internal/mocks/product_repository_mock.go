package mocks

import (
	"context"

	"product-services/internal/models"
	"product-services/internal/shared"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetProductByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) ListProducts(
	ctx context.Context,
	opts shared.ListOptions,
) (*models.ListProductsResult, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*models.ListProductsResult), args.Error(1)
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, category *models.Product) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, category *models.Product) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
