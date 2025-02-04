package types

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
	Arrdess string `json:"arrdess" validate:"required"`
}

type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity" validate:"gt=0"`
}
