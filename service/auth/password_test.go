package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {

	testCases := []struct {
		val    string
		errMsg string
	}{
		{val: "", errMsg: "expected hash to be not empty"},
		{val: "password", errMsg: "expected hash to be not empty"},
	}
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}
	for _, v := range testCases {
		if hash == v.val {
			t.Error(v.errMsg)
		}
	}
}

func TestComparePasswords(t *testing.T) {

	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if err := ComparePasswords([]byte(hash), []byte("password")); err != nil {
		t.Errorf("expected password to match hash")
	}
	
	if err := ComparePasswords([]byte(hash), []byte("passwordddddd")); err == nil {
		t.Errorf("expected password to not match hash")
	}
}
