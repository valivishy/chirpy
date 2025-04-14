package tests

import (
	"net/http"
	"testing"
)

func TestHandleHealthz(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	resp, err := http.Get(ts.BaseURL + "/api/healthz")
	if err != nil {
		t.Fatalf("failed to GET /api/healthz: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
		t.Errorf("expected content-type text/plain; charset=utf-8, got %s", contentType)
	}
}
