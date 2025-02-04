package product

import (
	"database/sql"

	"github.com/OlegB1/ecom/types"
	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts(p types.Pagination) ([]types.Product, error) {

	rows, err := s.db.Query("SELECT id, name, description, image, price, quantity, created_at FROM products LIMIT $1 OFFSET $2", p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) CreateProduct(p types.CreateProductPayload) (types.Product, error) {
	product := types.Product{}

	err := s.db.QueryRow(`
		INSERT INTO products (name,description,image,price,quantity) 
		VALUES($1,$2,$3,$4,$5)
		RETURNING id, name, description, image, price, quantity, created_at`,
		p.Name, p.Description, p.Image, p.Price, p.Quantity).Scan(
		&product.ID, &product.Name, &product.Description, &product.Image, &product.Price, &product.Quantity, &product.CreatedAt,
	)
	if err != nil {
		return types.Product{}, err
	}
	return product, nil
}

func (s *Store) GetProductsByIds(tx *sql.Tx, productIDs []int) ([]types.Product, error) {

	query := "SELECT id, name, description, image, price, quantity, created_at FROM products WHERE id = ANY($1)"

	rows, err := tx.Query(query, pq.Array(productIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) UpdateProducts(tx *sql.Tx, products []*types.Product) {
	ids := make([]int, len(products))
	quantities := make([]int, len(products))
	for i, u := range products {
		ids[i] = u.ID
		quantities[i] = u.Quantity
	}

	tx.Exec(`
		UPDATE products
		SET quantity = update_data.quantity
		FROM (SELECT unnest($1::int[]) AS id, unnest($2::int[]) AS quantity) AS update_data
		WHERE products.id = update_data.id
	`, pq.Array(ids), pq.Array(quantities))

}

func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return product, nil
}
