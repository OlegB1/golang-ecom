package types

import (
	"database/sql"
	"time"
)

type ProductStore interface {
	GetProducts(p Pagination) ([]Product, error)
	GetProductsByIds(*sql.Tx, []int) ([]Product, error)
	CreateProduct(CreateProductPayload) (Product, error)
	UpdateProducts(*sql.Tx, []*Product)
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       int       `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateProductPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
	Quantity    int64  `json:"quantity" validate:"required,gt=0"`
}
