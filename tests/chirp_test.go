package tests

import (
	"chirpy/models"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

const baseChirpsPath = "/api/chirps"

func TestHandleCreateChirp_Valid(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "chirper@example.com"
	password := "chirpy123"
	userId := createUser(t, ts, email, password)

	token, err := loginUser(t, ts, email, password)
	if err != nil {
		t.Fatal(err)
	}

	chirpBody := "Hello Chirpy!"
	body := `{"body":"` + chirpBody + `"}`

	var chirp models.ChirpDTO
	post(t, ts, baseChirpsPath, body, token, http.StatusCreated, &chirp)

	var chirps []models.ChirpDTO
	get(t, ts, baseChirpsPath, &chirps)

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
	defer closer(t)(ts.Server)

	token := createAndLoginUser(t, ts, "long@chirp.com", "chirpy456")

	longBody := strings.Repeat("a", 141)
	payload := fmt.Sprintf(`{"body":"%s"`, longBody)

	post(t, ts, baseChirpsPath, payload, token, http.StatusBadRequest, &models.ChirpDTO{})
}

func TestHandleCreateChirp_BannedWords(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	token := createAndLoginUser(t, ts, "filter@chirp.com", "chirpyfilter")

	body := `{"body":"This is a kerfuffle tweet"}`

	var chirp models.ChirpDTO
	post(t, ts, baseChirpsPath, body, token, http.StatusCreated, &chirp)

	if !strings.Contains(*chirp.Body, "****") {
		t.Errorf("expected banned word to be filtered: %s", *chirp.Body)
	}
}

func TestHandleCreateChirp_InvalidPayload(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	token := createAndLoginUser(t, ts, "filter_faulty@chirp.com", "chirpyfilter")

	payload := `{"body":false}`
	post(t, ts, baseChirpsPath, payload, token, http.StatusBadRequest, &models.ChirpDTO{})
}

func createAndLoginUser(
	t *testing.T, ts *TestServer, email, password string,
) string {
	createUser(t, ts, email, password)

	token, err := loginUser(t, ts, email, password)
	if err != nil {
		t.Fatal(err)
	}
	return token
}
