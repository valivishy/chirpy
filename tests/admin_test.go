package tests

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHandleAdminMetrics(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	resp, err := http.Get(ts.BaseURL + "/admin/metrics")
	if err != nil {
		t.Fatalf("failed to GET /admin/metrics: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("expected content-type text/html; charset=utf-8, got %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "Chirpy has been visited") {
		t.Errorf("response body did not contain expected text, got: %s", body)
	}
}

func TestHandleAdminReset(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	resp, err := http.Post(ts.BaseURL+"/admin/reset", "application/json", strings.NewReader(""))
	if err != nil {
		t.Fatalf("failed to POST /admin/reset: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
