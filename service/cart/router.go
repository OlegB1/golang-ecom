package cart

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/OlegB1/ecom/service/auth"
	"github.com/OlegB1/ecom/service/order"
	"github.com/OlegB1/ecom/service/product"
	"github.com/OlegB1/ecom/types"
	"github.com/OlegB1/ecom/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func RegisterHandler(db *sql.DB, router *mux.Router) {
	handler := NewHandler(order.NewStore(db), product.NewStore(db), db)

	router.HandleFunc("/cart/checkout", handler.HandleCheckout).Methods(http.MethodPost)
}

type Handler struct {
	orderStore   types.OrderStore
	productStore types.ProductStore
	DB           *sql.DB
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore, db *sql.DB) *Handler {
	return &Handler{orderStore: orderStore, productStore: productStore, DB: db}
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
