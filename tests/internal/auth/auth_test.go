package tests

import (
	"chirpy/internal/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"testing"
	"time"
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

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"
	expiresIn := time.Hour

	token, err := auth.MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	parsedUserID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("expected userID %s, got %s", userID, parsedUserID)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := auth.ValidateJWT("not.a.valid.token", "secret")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	token, err := auth.MakeJWT(userID, "rightsecret", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	_, err = auth.ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Bearer abc.def.ghi")

		token, err := auth.GetBearerToken(headers)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "abc.def.ghi" {
			t.Errorf("expected 'abc.def.ghi', got '%s'", token)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		headers := http.Header{}

		_, err := auth.GetBearerToken(headers)
		if err == nil {
			t.Fatal("expected error for missing header, got nil")
		}
	})

	t.Run("no Bearer prefix", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "abc.def.ghi")

		token, err := auth.GetBearerToken(headers)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "abc.def.ghi" {
			t.Errorf("expected token without prefix untouched, got '%s'", token)
		}
	})
}
