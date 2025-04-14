package tests

import (
	"chirpy/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"testing"
)

func TestHandleCreateUser(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	createUser(t, ts, "test@example.com", "superPassword123")
}

func TestHandleCreateUser_Duplicate(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	email := "test2@example.com"
	password := "superPassword1232"

	createUser(t, ts, email, password)
	payload := buildUserCreateOrLoginPayload(email, password)

	resp2, err := http.Post(ts.BaseURL+"/api/users", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("duplicate POST failed: %v", err)
	}
	defer Closer(t)(resp2.Body)

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request for duplicate user, got %d", resp2.StatusCode)
	}
}

func TestHandleCreateUser_InvalidPayload(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	invalidPayload := `{"email": 123, "password": true}`
	resp, err := http.Post(ts.BaseURL+"/api/users", "application/json", strings.NewReader(invalidPayload))
	if err != nil {
		t.Fatalf("POST with invalid payload failed: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request for invalid payload, got %d", resp.StatusCode)
	}
}

func buildUserCreateOrLoginPayload(email string, password string) string {
	return fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
}

func createUser(
	t *testing.T, ts *TestServer, email string, password string,
) uuid.UUID {
	payload := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)

	resp, err := http.Post(ts.BaseURL+"/api/users", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to POST /api/users: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	userResponse := models.UserDTO{}
	if err = decoder.Decode(&userResponse); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if !strings.Contains(userResponse.Email, email) {
		t.Errorf("expected email in response, got: %s", userResponse.Email)
	}

	return userResponse.ID
}
