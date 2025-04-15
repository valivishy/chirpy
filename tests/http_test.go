package tests

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

import (
	"net/http"
	"strings"
	"testing"
)

func get[T any](t *testing.T, ts *TestServer, url string, target *T) {
	t.Helper()

	resp, err := http.Get(ts.BaseURL + url)
	if err != nil {
		t.Fatalf("GET %s failed: %v", url, err)
	}
	defer closer(t)(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

func post[T any](
	t *testing.T, ts *TestServer, url string, body string, token string, expectedStatus int, target *T,
) {
	exec(t, ts, http.MethodPost, url, body, token, expectedStatus, target)
}

func put[T any](
	t *testing.T, ts *TestServer, url string, body string, token string, expectedStatus int, target *T,
) {
	exec(t, ts, http.MethodPut, url, body, token, expectedStatus, target)
}

func exec[T any](
	t *testing.T, ts *TestServer, method string, url string, body string, token string, expectedStatus int, target *T,
) {
	t.Helper()

	req, err := http.NewRequest(method, ts.BaseURL+url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s failed: %v", url, err)
	}
	defer closer(t)(resp.Body)

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected %d, got %d", expectedStatus, resp.StatusCode)
	}

	if expectedStatus != http.StatusOK && expectedStatus != http.StatusCreated {
		return
	}

	// We parse the response only on valid statuses
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}
