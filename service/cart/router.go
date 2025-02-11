package cart

import (
	"fmt"
	"net/http"

	"github.com/OlegB1/ecom/service/auth"
	"github.com/OlegB1/ecom/types"
	"github.com/OlegB1/ecom/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func RegisterHandler(orderRepository types.OrderRepository, productRepository types.ProductRepository, router *mux.Router) {
	handler := NewHandler(orderRepository, productRepository)

	router.HandleFunc("/cart/checkout", handler.HandleCheckout).Methods(http.MethodPost)
}

type Handler struct {
	orderRepository   types.OrderRepository
	productRepository types.ProductRepository
}

func NewHandler(orderStore types.OrderRepository, productStore types.ProductRepository) *Handler {
	return &Handler{orderRepository: orderStore, productRepository: productStore}
}

func (h *Handler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	var cart types.CartCheckoutPayload
	if err := utils.ParseJson(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())
	if userID == -1 {
		fmt.Println("wrong user id from token")
		utils.PermissionDanied(w)
	}
	order, err := h.CreateOrder(cart, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, order)

}
