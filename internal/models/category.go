package models

import (
	"time"

	"github.com/google/uuid"
)

// Common types
type Pagination struct {
	NextCursor time.Time
	HasMore    bool
}

type TimeStamps struct {
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// Category models
type Category struct {
	ID          uuid.UUID `json:"id"          db:"id"`
	Name        string    `json:"name"        db:"name"`
	Description string    `json:"description" db:"description"`
	TimeStamps
}

type ListCategoriesResult struct {
	Categories []*Category
	Pagination
}

type CategoryRequest struct {
	Name        string `json:"name"        validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

// Product models
type Product struct {
	ID          uuid.UUID `json:"id"          db:"id"`
	Name        string    `json:"name"        db:"name"`
	Description string    `json:"description" db:"description"`
	ImageURL    string    `json:"imageUrl"    db:"image_url"`
	CategoryID  uuid.UUID `json:"categoryID"  db:"category_id"`
	Price       float64   `json:"price"       db:"price"`
	Quantity    int       `json:"quantity"    db:"quantity"`
	TimeStamps
}

type ListProductsResult struct {
	Products []*Product
	Pagination
}

type ProductRequest struct {
	Name        string    `json:"name"        validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"omitempty,max=255"`
	ImageURL    string    `json:"imageUrl"    validate:"omitempty,max=255"`
	CategoryID  uuid.UUID `json:"categoryID"  validate:"required"`
	Price       float64   `json:"price"       validate:"required"`
	Quantity    int       `json:"quantity"    validate:"required"`
}
