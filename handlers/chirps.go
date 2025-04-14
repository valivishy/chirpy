package handlers

import (
	"chirpy/config"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/models"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

func HandleCreateChirp(api *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userId, err := auth.ValidateJWT(token, api.Secret)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(r.Body)
		createChirpRequest := models.CreateChirpRequest{}
		if err = decoder.Decode(&createChirpRequest); err != nil {
			respondWithError(w, "Invalid payload", http.StatusBadRequest)
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
				UserID: userId,
			},
		)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		printJsonResponse(w, mapChirp(chirp), http.StatusCreated)
	}
}

func HandleListChirps(api *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		chirps, err := api.Queries.ListChirps(r.Context())
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		var results []models.ChirpDTO
		for _, chirp := range chirps {
			results = append(results, mapChirp(chirp))
		}

		printJsonResponse(w, results, http.StatusOK)
	}
}

func HandleGetChirp(api *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
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

func mapChirp(chirp database.Chirp) models.ChirpDTO {
	return models.ChirpDTO{
		ID:        &chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      &chirp.Body,
		UserID:    &chirp.UserID,
	}
}
