package tests

import (
	"chirpy/models"
	"fmt"
	_ "github.com/lib/pq"
)

import (
	"net/http"
	"testing"
)

func buildUserRequestPayload(email string, password string) string {
	return fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
}

func loginUser(t *testing.T, ts *TestServer, email, password string) (*models.UserDTO, error) {
	loginPayload := buildUserRequestPayload(email, password)

	var user models.UserDTO
	post(t, ts, "/api/login", loginPayload, "", http.StatusOK, &user)

	return &user, nil
}
