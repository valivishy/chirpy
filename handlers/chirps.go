package handlers

import (
	"chirpy/config"
	"chirpy/internal/database"
	"chirpy/models"
	"github.com/google/uuid"
	"net/http"
)

func HandleCreateChirp(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := validateJWT(configuration, r.Header)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		createChirpRequest, err := decodeRequestPayload[models.CreateChirpRequest](r)
		if err != nil {
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

		chirp, err := configuration.Queries.CreateChirp(
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

func HandleDeleteChirp(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chirp, failed := getChirp(w, r, configuration)
		if failed {
			return
		}

		if err := configuration.Queries.DeleteChirp(r.Context(), chirp.ID); err != nil {
			respondWithError(w, err.Error(), http.StatusNotFound)
			return
		}

		printJsonResponse(w, "", http.StatusNoContent)
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

func getChirp(w http.ResponseWriter, r *http.Request, configuration *config.Configuration) (database.Chirp, bool) {
	userId, err := validateJWT(configuration, r.Header)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusUnauthorized)
		return database.Chirp{}, true
	}

	value := r.PathValue("chirpID")
	if value == "" {
		respondWithError(w, "ChirpID path variable is mandatory", http.StatusBadRequest)
	}

	chirp, err := configuration.Queries.GetChirp(r.Context(), uuid.MustParse(value))
	if err != nil {
		respondWithError(w, err.Error(), http.StatusNotFound)
		return database.Chirp{}, true
	}

	if chirp.UserID != userId {
		respondWithError(w, "", http.StatusForbidden)
		return database.Chirp{}, true
	}

	if len(chirp.Body) == 0 {
		respondWithError(w, "", http.StatusNotFound)
		return database.Chirp{}, true
	}

	return chirp, false
}
