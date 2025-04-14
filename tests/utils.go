package tests

import (
	"chirpy/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type TestServer struct {
	Server  *http.Server
	BaseURL string
}

func closer(t *testing.T) func(io.Closer) {
	return func(c io.Closer) {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func buildUserCreateOrLoginPayload(email string, password string) string {
	return fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)
}

func loginUser(t *testing.T, ts *TestServer, email, password string) (string, error) {
	loginPayload := buildUserCreateOrLoginPayload(email, password)
	resp, err := http.Post(ts.BaseURL+"/api/login", "application/json", strings.NewReader(loginPayload))
	if err != nil {
		return "", err
	}
	defer closer(t)(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	var user models.UserDTO
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	return user.Token, nil
}
