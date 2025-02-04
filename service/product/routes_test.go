package product

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OlegB1/ecom/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestProductServiceHandlers(t *testing.T) {

	t.Run("Should failed get products", func(t *testing.T) {
		mockPoductStore := new(MockProductStore)
		handler := NewHandler(mockPoductStore)

		mockPoductStore.On("GetProducts", types.Pagination{Offset: 4, Limit: 4}).Return([]types.Product{}, errors.New("database error"))

		req, err := http.NewRequest(http.MethodGet, "/producs?offset=4&limit=4", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/producs", handler.HandleGetProducts)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("unexpected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		mockPoductStore.AssertExpectations(t)

	})

	t.Run("Should get products", func(t *testing.T) {
		mockPoductStore := new(MockProductStore)
		handler := NewHandler(mockPoductStore)

		mockPoductStore.On("GetProducts", types.Pagination{Offset: 4, Limit: 4}).Return([]types.Product{
			{
				ID:          0,
				Name:        "Product 1",
				Description: "Description 1",
				Image:       "image_url",
				Price:       100,
				Quantity:    10,
				CreatedAt:   time.Now(),
			},
		}, nil)

		req, err := http.NewRequest(http.MethodGet, "/producs?offset=4&limit=4", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/producs", handler.HandleGetProducts)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("unexpected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockPoductStore.AssertExpectations(t)

	})

	t.Run("Should faile create product on validation", func(t *testing.T) {
		product := types.CreateProductPayload{
			Name:        "name",
			Description: "des",
			Image:       "img",
			Price:       -1,
			Quantity:    -1,
		}
		marshaled, _ := json.Marshal(product)

		mockPoductStore := new(MockProductStore)
		handler := NewHandler(mockPoductStore)

		req, err := http.NewRequest(http.MethodPost, "/producs", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/producs", handler.HandleCreateProduct)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("unexpected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockPoductStore.AssertExpectations(t)

	})

	t.Run("Should failed create product", func(t *testing.T) {
		product := types.CreateProductPayload{
			Name:        "name",
			Description: "des",
			Image:       "img",
			Price:       1,
			Quantity:    1,
		}
		marshaled, _ := json.Marshal(product)

		mockPoductStore := new(MockProductStore)
		handler := NewHandler(mockPoductStore)

		mockPoductStore.On("CreateProduct", product).Return(types.Product{}, errors.New("database error"))

		req, err := http.NewRequest(http.MethodPost, "/producs", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/producs", handler.HandleCreateProduct)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("unexpected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockPoductStore.AssertExpectations(t)
	})

	t.Run("Should create product", func(t *testing.T) {
		product := types.CreateProductPayload{
			Name:        "name",
			Description: "des",
			Image:       "img",
			Price:       1,
			Quantity:    1,
		}
		marshaled, _ := json.Marshal(product)

		mockPoductStore := new(MockProductStore)
		handler := NewHandler(mockPoductStore)

		mockPoductStore.On("CreateProduct", product).Return(types.Product{
			ID:          1,
			Name:        product.Image,
			Description: product.Description,
			Image:       product.Image,
			Price:       int(product.Price),
			Quantity:    int(product.Quantity),
			CreatedAt:   time.Time{},
		}, nil)

		req, err := http.NewRequest(http.MethodPost, "/producs", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/producs", handler.HandleCreateProduct)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("unexpected status code %d, got %d", http.StatusOK, rr.Code)
		}

		mockPoductStore.AssertExpectations(t)
	})
}

// Mock struct definition
type MockProductStore struct {
	mock.Mock
}

func (s *MockProductStore) GetProducts(p types.Pagination) ([]types.Product, error) {
	args := s.Called(p)
	return args.Get(0).([]types.Product), args.Error(1)
}

func (s *MockProductStore) GetProductsByIds(tx *sql.Tx, ids []int) ([]types.Product, error) {
	args := s.Called(tx, ids)
	return args.Get(0).([]types.Product), args.Error(1)
}

func (s *MockProductStore) CreateProduct(payload types.CreateProductPayload) (types.Product, error) {
	args := s.Called(payload)
	return args.Get(0).(types.Product), args.Error(1)
}

func (s *MockProductStore) UpdateProducts(tx *sql.Tx, products []*types.Product) {
	s.Called(tx, products)
}
