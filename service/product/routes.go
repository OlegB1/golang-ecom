package product

import (
	"fmt"
	"net/http"

	"github.com/OlegB1/ecom/types"
	"github.com/OlegB1/ecom/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func RegisterHandler(repository types.ProductRepository, router *mux.Router) {
	handler := NewHandler(repository)

	router.HandleFunc("/products", handler.HandleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", handler.HandleCreateProduct).Methods(http.MethodPost)
}

type Handler struct {
	repository types.ProductRepository
}

func NewHandler(store types.ProductRepository) Handler {
	return Handler{repository: store}
}

func (h *Handler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	pagination := utils.GetPagination(r)
	products, err := h.repository.GetProducts(pagination)

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
	product, err := h.repository.CreateProduct(p)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

}
