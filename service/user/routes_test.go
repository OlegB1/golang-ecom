package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OlegB1/ecom/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)
	t.Run("Should register failed user payload is invalid", func(t *testing.T) {
		peyload := types.RegisterUserPayload{
			FirstName: "ff",
			LastName:  "ll",
			Email:     "invalid",
			Password:  "123",
		}

		marshaled, _ := json.Marshal(peyload)
		req, err := http.NewRequest(http.MethodPatch, "/register", bytes.NewBuffer(marshaled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("nexpecteed status code %d , got %d", http.StatusBadRequest, rr.Code)
		}
	})
	t.Run("Shoult register success", func(t *testing.T) {
		peyload := types.RegisterUserPayload{
			FirstName: "ff",
			LastName:  "ll",
			Email:     "valid@mail.com",
			Password:  "123",
		}

		marshaled, _ := json.Marshal(peyload)
		req, err := http.NewRequest(http.MethodPatch, "/register", bytes.NewBuffer(marshaled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("nexpecteed status code %d , got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}
func (m *mockUserStore) GetUserById(id int) (*types.User, error) {
	return nil, nil
}
func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}
