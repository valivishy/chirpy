package handlers

import (
	"chirpy/config"
	"chirpy/internal/database"
	"chirpy/models"
	"github.com/google/uuid"
	"net/http"
	"sort"
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
		userId := r.URL.Query().Get("author_id")
		sortParam := getSortParam(r)

		var chirps []database.Chirp
		var err error
		if userId == "" {
			chirps, err = api.Queries.ListAllChirps(r.Context())
		} else {
			userID, err2 := uuid.Parse(userId)
			if err2 != nil {
				respondWithError(w, err2.Error(), http.StatusBadRequest)
				return
			}

			chirps, err = api.Queries.ListChirpsByUser(r.Context(), userID)
		}
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		var results []models.ChirpDTO
		for _, chirp := range chirps {
			results = append(results, mapChirp(chirp))
		}

		sort.Slice(results, func(i, j int) bool {
			if sortParam == "desc" {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			}
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})

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

func getSortParam(r *http.Request) string {
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "asc" && sortParam != "desc" {
		sortParam = "asc"
	}
	return sortParam
}
