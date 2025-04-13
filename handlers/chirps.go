package handlers

import (
	"chirpy/config"
	"chirpy/internal/database"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type ChirpDTO struct {
	Error     *string    `json:"error"`
	Valid     *bool      `json:"valid"`
	ID        *uuid.UUID `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Body      *string    `json:"body"`
	UserID    *uuid.UUID `json:"user_id"`
}

func HandleCreateChirp(api *config.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		createChirpRequest := CreateChirpRequest{}
		err := decoder.Decode(&createChirpRequest)
		if err != nil {
			respondWithError(w, "Something went wrong", http.StatusBadRequest)
			return
		}

		if len(createChirpRequest.Body) > 140 {
			respondWithError(w, "Chirp is too long", http.StatusBadRequest)
			return
		}

		text := createChirpRequest.Body
		for _, word := range []string{"kerfuffle ", "sharbert ", "fornax "} {
			text = replaceInsensitive(text, word, "**** ")
		}

		chirp, err := api.Queries.CreateChirp(
			r.Context(),
			database.CreateChirpParams{
				Body:   text,
				UserID: createChirpRequest.UserId,
			},
		)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		printJsonResponse(w, mapChirp(chirp), http.StatusCreated)
	}
}

func HandleListChirps(api *config.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		chirps, err := api.Queries.ListChirps(r.Context())
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		var results []ChirpDTO
		for _, chirp := range chirps {
			results = append(results, mapChirp(chirp))
		}

		printJsonResponse(w, results, http.StatusOK)
	}
}

func HandleGetChirp(api *config.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		value := r.PathValue("chirpID")
		if value == "" {
			respondWithError(w, "ChirpID path variable is mandatory", http.StatusBadRequest)
		}

		chirp, err := api.Queries.GetChirp(r.Context(), uuid.MustParse(value))
		if err != nil {
			respondWithError(w, err.Error(), http.StatusNotFound)
			return
		}

		printJsonResponse(w, mapChirp(chirp), http.StatusOK)
	}
}

func mapChirp(chirp database.Chirp) ChirpDTO {
	return ChirpDTO{
		ID:        &chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      &chirp.Body,
		UserID:    &chirp.UserID,
	}
}
