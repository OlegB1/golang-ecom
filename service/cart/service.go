package cart

import (
	"database/sql"
	"fmt"

	"github.com/OlegB1/ecom/types"
)

func (h *Handler) GetCartItemsIDs(items []types.CartItem) []int {
	productIds := []int{}
	for _, v := range items {
		productIds = append(productIds, v.ProductID)
	}
	return productIds
}

func (h *Handler) CreateOrder(cart types.CartCheckoutPayload, userID int) (types.Order, error) {
	if len(cart.Items) == 0 {
		return types.Order{}, fmt.Errorf("cart is empy")
	}

	tx, err := h.DB.Begin()
	if err != nil {
		return types.Order{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// get products
	productIDs := h.GetCartItemsIDs(cart.Items)
	products, _ := h.productStore.GetProductsByIds(tx, productIDs)

	productMap := make(map[int]*types.Product)

	for _, v := range products {
		productMap[v.ID] = &v
	}

	// check if products in stock
	if err := h.CheckProductsInStock(tx, cart.Items, productMap); err != nil {
		return types.Order{}, err
	}

	// calc total price
	totalPrice := 0
	for _, v := range cart.Items {
		product := productMap[v.ProductID]
		totalPrice += v.Quantity * product.Price
		//reduce quantity of product
		product.Quantity = product.Quantity - v.Quantity
	}

	// update products in db
	productsToUpdate := make([]*types.Product, len(cart.Items))
	for i, v := range cart.Items {
		productsToUpdate[i] = productMap[v.ProductID]
	}
	h.productStore.UpdateProducts(tx, productsToUpdate)

	// create order
	order := h.orderStore.CreateOrder(tx, types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: cart.Arrdess,
	})
	if err != nil {
		return types.Order{}, err
	}

	// create order-items
	orderItems := make([]types.OrederItem, len(productsToUpdate))
	for i, item := range cart.Items {
		orderItems[i] = types.OrederItem{
			OrderID:   order.ID,
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		}
	}
	h.orderStore.CreateOrderItems(tx, orderItems)

	return order, nil
}

func (h *Handler) CheckProductsInStock(tx *sql.Tx, items []types.CartItem, productMap map[int]*types.Product) error {
	for _, v := range items {
		product, ok := productMap[v.ProductID]
		if !ok {
			return fmt.Errorf("product with id %d is not available in store", v.ProductID)
		}

		if product.Quantity < v.Quantity {
			return fmt.Errorf("product with id %d is not available in the requested quantity", v.ProductID)
		}

	}
	return nil
}
