package types

import (
	"database/sql"
	"time"
)

type OrderRepository interface {
	CreateOrder(*sql.Tx, Order) Order
	CreateOrderItems(*sql.Tx, []OrederItem) error
	GetTransaction() (*sql.Tx, error)
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Total     int       `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrederItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderId"`
	ProductId int       `json:"productOd"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}
