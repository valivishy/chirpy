package tests

import (
	"chirpy/models"
	"net/http"
	"testing"
)

func TestHandleRefresh_Success(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "valid+1@example.com"
	password := "correctPassword"
	createUser(t, ts, email, password)

	user, err := loginUser(t, ts, email, password)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	var response models.RefreshTokenResponse
	execPost(t, ts, "/api/refresh", "", user.RefreshToken, http.StatusOK, &response)

	if len(response.Token) == 0 {
		t.Fatalf("Refresh failed: %v", err)
	}
}

func TestHandleRevoke_Success(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "valid+12@example.com"
	password := "correctPassword"
	createUser(t, ts, email, password)

	user, err := loginUser(t, ts, email, password)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	var response models.RefreshTokenResponse
	execPost(t, ts, "/api/revoke", "", user.RefreshToken, http.StatusNoContent, &response)

	execPost(t, ts, "/api/refresh", "", user.RefreshToken, http.StatusUnauthorized, &response)
}
