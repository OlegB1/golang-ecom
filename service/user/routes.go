package user

import (
	"fmt"
	"net/http"

	"github.com/OlegB1/ecom/service/auth"
	"github.com/OlegB1/ecom/types"
	"github.com/OlegB1/ecom/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func RegisterHandler(repository types.UserRepository, router *mux.Router) {
	handler := NewHandler(repository)
	router.HandleFunc("/login", handler.handleLogin).Methods("POST")
	router.HandleFunc("/register", handler.handleRegister).Methods("POST")

}

type Handler struct {
	repository types.UserRepository
}

func NewHandler(store types.UserRepository) Handler {
	return Handler{repository: store}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	u, err := h.repository.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found email %s", payload.Email))
		return
	}

	if err = auth.ComparePasswords([]byte(u.Password), []byte(payload.Password)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("wrong password"))
		return
	}

	token, err := auth.CreateJWT(u.ID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	if _, err := h.repository.GetUserByEmail(payload.Email); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return

	}

	hashedPass, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	err = h.repository.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPass,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
