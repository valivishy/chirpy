package tests

import (
	"net/http"
	"strings"
	"testing"
)

func TestHandleLogin_Success(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "valid@example.com"
	password := "correctPassword"
	createUser(t, ts, email, password)

	loginPayload := buildUserCreateOrLoginPayload(email, password)
	resp, err := http.Post(ts.BaseURL+"/api/login", "application/json", strings.NewReader(loginPayload))
	if err != nil {
		t.Fatalf("failed to POST /api/login: %v", err)
	}
	defer closer(t)(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestHandleLogin_UserNotFound(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	loginPayload := buildUserCreateOrLoginPayload("nouser@example.com", "irrelevant")
	resp, err := http.Post(ts.BaseURL+"/api/login", "application/json", strings.NewReader(loginPayload))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer closer(t)(resp.Body)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for missing user, got %d", resp.StatusCode)
	}
}

func TestHandleLogin_WrongPassword(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "user@example.com"
	correctPassword := "goodPass"
	wrongPassword := "badPass"
	createUser(t, ts, email, correctPassword)

	loginPayload := buildUserCreateOrLoginPayload(email, wrongPassword)
	resp, err := http.Post(ts.BaseURL+"/api/login", "application/json", strings.NewReader(loginPayload))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer closer(t)(resp.Body)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for wrong password, got %d", resp.StatusCode)
	}
}

func TestHandleLogin_InvalidEmailFormat(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	payload := buildUserCreateOrLoginPayload("1234", "irrelevant")
	resp, err := http.Post(ts.BaseURL+"/api/login", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer closer(t)(resp.Body)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 400 Bad Request for invalid payload, got %d", resp.StatusCode)
	}
}
