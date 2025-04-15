package tests

import (
	"chirpy/models"
	"net/http"
	"testing"
)

func TestHandleLogin_Success(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "valid@example.com"
	password := "correctPassword"
	createUser(t, ts, email, password)

	_, err := loginUser(t, ts, email, password)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
}

func TestHandleLogin_UserNotFound(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	loginPayload := buildUserCreateOrLoginPayload("nouser@example.com", "irrelevant")
	post(t, ts, "/api/login", loginPayload, "", http.StatusUnauthorized, &models.UserDTO{})
}

func TestHandleLogin_WrongPassword(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "user@example.com"
	correctPassword := "goodPass"
	wrongPassword := "badPass"
	createUser(t, ts, email, correctPassword)

	loginPayload := buildUserCreateOrLoginPayload(email, wrongPassword)
	post(t, ts, "/api/login", loginPayload, "", http.StatusUnauthorized, &models.UserDTO{})
}

func TestHandleLogin_InvalidEmailFormat(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	payload := buildUserCreateOrLoginPayload("1234", "irrelevant")
	post(t, ts, "/api/login", payload, "", http.StatusUnauthorized, &models.UserDTO{})
}
