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

	user := createAndLoginUser(t, ts, "chirper@example.com", "chirpy123")

	createChirp(t, ts, &user)
}

func TestHandleCreateChirp_TooLong(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	user := createAndLoginUser(t, ts, "long@chirp.com", "chirpy456")

	longBody := strings.Repeat("a", 141)
	payload := fmt.Sprintf(`{"body":"%s"`, longBody)

	execPost(t, ts, baseChirpsPath, payload, user.Token, "Bearer", http.StatusBadRequest, &models.ChirpDTO{})
}

func TestHandleCreateChirp_BannedWords(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	user := createAndLoginUser(t, ts, "filter@chirp.com", "chirpyfilter")

	body := `{"body":"This is a kerfuffle tweet"}`

	var chirp models.ChirpDTO
	execPost(t, ts, baseChirpsPath, body, user.Token, "Bearer", http.StatusCreated, &chirp)

	if !strings.Contains(*chirp.Body, "****") {
		t.Errorf("expected banned word to be filtered: %s", *chirp.Body)
	}
}

func TestHandleCreateChirp_InvalidPayload(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	user := createAndLoginUser(t, ts, "filter_faulty@chirp.com", "chirpyfilter")

	payload := `{"body":false}`
	execPost(t, ts, baseChirpsPath, payload, user.Token, "Bearer", http.StatusBadRequest, &models.ChirpDTO{})
}

func TestHandleDeleteChirp_Valid(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "chirper_delete_chirp@example.com"
	password := "chirpy123"
	user := createAndLoginUser(t, ts, email, password)

	chirp := createChirp(t, ts, &user)

	execDelete(t, ts, "/api/chirps/"+chirp.ID.String(), user.Token, http.StatusNoContent)

	var newChirp models.ChirpDTO
	get(t, ts, "/api/chirps/"+chirp.ID.String(), user.Token, http.StatusNotFound, &newChirp)
	if newChirp.Body != nil {
		t.Fatal("chirp found")
	}
}

func TestHandleDeleteChirpByAnotherUser_Fails(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "chirper_delete_fails@example.com"
	password := "chirpy123"
	user := createAndLoginUser(t, ts, email, password)
	chirp := createChirp(t, ts, &user)

	user2 := createAndLoginUser(t, ts, "chirper_delete_fails2@example.com", "chirpy12322")

	execDelete(t, ts, "/api/chirps/"+chirp.ID.String(), user2.Token, http.StatusForbidden)
}

func TestHandleListChirps_SortingWorks(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	user := createAndLoginUser(t, ts, "chirper_list_chirps_sort@example.com", "chirpy123")

	createChirp(t, ts, &user)
	createChirp(t, ts, &user)
	createChirp(t, ts, &user)
	createChirp(t, ts, &user)
	createChirp(t, ts, &user)

	chirps := listChirps(t, ts, &user, "")
	for i := 0; i < len(chirps)-1; i++ {
		if chirps[i].CreatedAt.After(chirps[i+1].CreatedAt) {
			t.Fatal("Incorrect sorting")
		}
	}

	chirps = listChirps(t, ts, &user, "desc")
	for i := 0; i < len(chirps)-1; i++ {
		if chirps[i].CreatedAt.Before(chirps[i+1].CreatedAt) {
			t.Fatal("Incorrect sorting")
		}
	}
}

func createAndLoginUser(
	t *testing.T, ts *TestServer, email, password string,
) models.UserDTO {
	createUser(t, ts, email, password)

	user, err := loginUser(t, ts, email, password)
	if err != nil {
		t.Fatal(err)
	}
	return *user
}

func createChirp(t *testing.T, ts *TestServer, user *models.UserDTO) models.ChirpDTO {
	chirpBody := "Hello Chirpy!"
	body := `{"body":"` + chirpBody + `"}`

	var chirp models.ChirpDTO
	execPost(t, ts, baseChirpsPath, body, user.Token, "Bearer", http.StatusCreated, &chirp)

	for _, listedChirp := range listChirps(t, ts, user, "") {
		if user.ID == *listedChirp.UserID && *listedChirp.Body == chirpBody {
			chirp = listedChirp
			break
		}
	}
	if chirp.ID == nil {
		t.Errorf("chirp not created")
		return models.ChirpDTO{}
	}

	get(t, ts, "/api/chirps/"+chirp.ID.String(), user.Token, http.StatusOK, &chirp)
	if chirp.UserID == nil {
		t.Errorf("chirp was not found")
		return models.ChirpDTO{}
	}

	return chirp
}

func listChirps(t *testing.T, ts *TestServer, user *models.UserDTO, sort string) []models.ChirpDTO {
	var chirps []models.ChirpDTO
	get(t, ts, baseChirpsPath+"?sort="+sort, user.Token, http.StatusOK, &chirps)
	return chirps
}
