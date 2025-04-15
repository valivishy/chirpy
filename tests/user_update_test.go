package tests

import (
	"net/http"
	"testing"
)

func TestHandleUpdate(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "test-update@example.com"
	password := "superPassword123"
	createUser(t, ts, email, password)

	user, err := loginUser(t, ts, email, password)
	if err != nil {
		t.Fatal(err)
	}

	newEmail := "updating_test-update@example.com"
	newPassword := "updating_superPassword123"

	execPut(t, ts, "/api/users", buildUserRequestPayload(newEmail, newPassword), user.Token, http.StatusOK, &user)

	execPost(t, ts, "/api/login", buildUserRequestPayload(email, password), "", "Bearer", http.StatusUnauthorized, &user)

	user, err = loginUser(t, ts, newEmail, newPassword)
	if err != nil {
		t.Fatal(err)
	}
}
