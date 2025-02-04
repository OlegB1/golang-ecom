package auth

import (
	"testing"
)

func TestCreateJWT(t *testing.T) {

	jwt, err := CreateJWT(1)

	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if jwt == "" {
		t.Error("Expected toketn to be not empty")
	}
}
