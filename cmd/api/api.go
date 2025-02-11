package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/OlegB1/ecom/service/auth"
	"github.com/OlegB1/ecom/service/cart"
	"github.com/OlegB1/ecom/service/order"
	"github.com/OlegB1/ecom/service/product"
	"github.com/OlegB1/ecom/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	router.Use(auth.JWTMiddleware(user.NewRepository(s.db)))

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userRepository := user.NewRepository(s.db)
	user.RegisterHandler(userRepository, subrouter)

	productRepository := product.NewRepository(s.db)
	product.RegisterHandler(productRepository, subrouter)

	orderRepository := order.NewRepository(s.db)
	cart.RegisterHandler(orderRepository, productRepository, subrouter)

	log.Println("Listening on ", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}
