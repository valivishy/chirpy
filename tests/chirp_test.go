package tests

import (
	"chirpy/models"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHandleCreateChirp_Valid(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	email := "chirper@example.com"
	password := "chirpy123"
	userId := createUser(t, ts, email, password)
	chirpBody := "Hello Chirpy!"
	body := `{"body":"` + chirpBody + `", "user_id":"` + userId.String() + `"}`

	resp, err := http.Post(ts.BaseURL+"/api/chirps", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST /api/chirps failed: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", resp.StatusCode)
	}

	var chirps []models.ChirpDTO
	get(t, ts, "/api/chirps", &chirps)

	var chirp models.ChirpDTO
	for _, listedChirp := range chirps {
		if userId == *listedChirp.UserID && *listedChirp.Body == chirpBody {
			chirp = listedChirp
			break
		}
	}
	if chirp.ID == nil {
		t.Errorf("chirp not created")
		return
	}

	get(t, ts, "/api/chirps/"+chirp.ID.String(), &chirp)
	if chirp.UserID == nil {
		t.Errorf("chirp was not found")
		return
	}
}

func TestHandleCreateChirp_TooLong(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	email := "long@chirp.com"
	password := "chirpy456"
	userId := createUser(t, ts, email, password)

	longBody := strings.Repeat("a", 141)
	payload := fmt.Sprintf(`{"body":"%s", "user_id":"%s"}`, longBody, userId)

	resp, err := http.Post(ts.BaseURL+"/api/chirps", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("POST /api/chirps failed: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for too long chirp, got %d", resp.StatusCode)
	}
}

func TestHandleCreateChirp_BannedWords(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	email := "filter@chirp.com"
	password := "chirpyfilter"
	userId := createUser(t, ts, email, password)

	body := `{"body":"This is a kerfuffle tweet", "user_id":"` + userId.String() + `"}`

	resp, err := http.Post(ts.BaseURL+"/api/chirps", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", resp.StatusCode)
	}

	content, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(content), "****") {
		t.Errorf("expected banned word to be filtered: %s", content)
	}
}

func TestHandleCreateChirp_InvalidPayload(t *testing.T) {
	ts := Start(t)
	defer Closer(t)(ts.Server)

	payload := `{"body":false, "user_id":123}`
	resp, err := http.Post(ts.BaseURL+"/api/chirps", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer Closer(t)(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid payload, got %d", resp.StatusCode)
	}
}
