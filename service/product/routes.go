package product

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/OlegB1/ecom/types"
	"github.com/OlegB1/ecom/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func RegisterHandler(db *sql.DB, router *mux.Router) {
	handler := NewHandler(&Store{db: db})

	router.HandleFunc("/producs", handler.HandleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/producs", handler.HandleCreateProduct).Methods(http.MethodPost)
}

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	pagination := utils.GetPagination(r)
	products, err := h.store.GetProducts(pagination)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, products); err != nil {
		fmt.Println(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
}

func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p types.CreateProductPayload

	if err := utils.ParseJson(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(p); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errs))
			return
		}
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	product, err := h.store.CreateProduct(p)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

}
