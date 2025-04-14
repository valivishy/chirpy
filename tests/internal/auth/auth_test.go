package tests

import (
	"chirpy/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword_Success(t *testing.T) {
	password := "supersecurepassword"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(hash) == 0 {
		t.Error("expected a non-empty hash")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Errorf("expected hash to match password, got error: %v", err)
	}
}

func TestCheckPasswordHash_Success(t *testing.T) {
	password := "mysecret"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := auth.CheckPasswordHash(string(hash), password)
	if err != nil {
		t.Errorf("expected password to match hash, got error: %v", err)
	}
}

func TestCheckPasswordHash_Failure(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	err := auth.CheckPasswordHash(string(hash), "wrongpassword")
	if err == nil {
		t.Error("expected error for mismatched password, got nil")
	}
}
