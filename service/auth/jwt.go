package auth

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/OlegB1/ecom/config"
	"github.com/OlegB1/ecom/types"
	"github.com/OlegB1/ecom/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type contextKey string

const UserKey contextKey = "userID"

var SKIP_JWT_URLS = []string{"login", "register"}

func CreateJWT(userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWT_EXPIRATION_SECONDS)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":    strconv.Itoa(userId),
			"expiredAt": time.Now().Add(expiration).Unix(),
		})

	tokenString, err := token.SignedString([]byte(config.Envs.JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTMiddleware(store types.UserStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропуск JWT перевірки для певних URL
			for _, v := range SKIP_JWT_URLS {
				if strings.Contains(r.URL.Path, v) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Отримання токена з заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("missing Authorization header")
				utils.PermissionDanied(w)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Структура claims
			type MyCustomClaims struct {
				UserID string `json:"userId"`
				jwt.RegisteredClaims
			}

			token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Envs.JWT_SECRET_KEY), nil
			})

			if err != nil || !token.Valid {
				log.Printf("failed to validate token: %v", err)
				utils.PermissionDanied(w)
				return
			}

			claims, ok := token.Claims.(*MyCustomClaims)
			if !ok {
				log.Println("unknown claims type")
				utils.PermissionDanied(w)
				return
			}

			userID, err := strconv.Atoi(claims.UserID)
			if err != nil {
				log.Println("invalid user ID in token")
				utils.PermissionDanied(w)
				return
			}

			user, err := store.GetUserById(userID)
			if err != nil {
				log.Printf("failed to get user by ID: %v", err)
				utils.PermissionDanied(w)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user.ID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}
