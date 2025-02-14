package order

import (
	"database/sql"

	"github.com/OlegB1/ecom/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (s *Repository) CreateOrder(tx *sql.Tx, o types.Order) types.Order {
	order := types.Order{}

	tx.QueryRow(`
		INSERT INTO orders (user_id, total, status, address)
		VALUES($1,$2,$3,$4)
		RETURNING id, user_id, total, status, address, created_at`,
		o.UserID, o.Total, o.Status, o.Address).Scan(
		&order.ID, &order.UserID, &order.Total, &order.Status, &order.Address, &order.CreatedAt,
	)

	return order
}

func (s *Repository) GetTransaction() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *Repository) CreateOrderItems(tx *sql.Tx, items []types.OrederItem) error {
	for _, item := range items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, price) 
			VALUES($1, $2, $3, $4)`,
			item.OrderID, item.ProductId, item.Quantity, item.Price)

		if err != nil {
			return err
		}
	}

	return nil

}
