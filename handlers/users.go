package handlers

import (
	"chirpy/config"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/mappers"
	"chirpy/models"
	"encoding/json"
	"net/http"
)

const somethingWentWrong = "Something went wrong"

func HandleCreate(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestBody := models.UserRequest{}
		err := decoder.Decode(&requestBody)

		if err != nil {
			respondWithError(w, somethingWentWrong, http.StatusBadRequest)
			return
		}

		password, err := auth.HashPassword(requestBody.Password)
		if err != nil {
			respondWithError(w, somethingWentWrong, http.StatusInternalServerError)
			return
		}

		user, err := configuration.Queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          requestBody.Email,
			HashedPassword: password,
		})
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
		}

		printJsonResponse(w, mappers.MapUser(user, "", ""), http.StatusCreated)
	}
}

func HandleUpdate(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestBody := models.UserRequest{}
		if err := decoder.Decode(&requestBody); err != nil {
			respondWithError(w, somethingWentWrong, http.StatusBadRequest)
			return
		}

		userId, err := validateJWT(configuration, r.Header)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		password, err := auth.HashPassword(requestBody.Password)
		if err != nil {
			respondWithError(w, somethingWentWrong, http.StatusInternalServerError)
			return
		}

		err = configuration.Queries.UpdateUser(r.Context(), database.UpdateUserParams{
			Email:          requestBody.Email,
			HashedPassword: password,
			ID:             userId,
		})
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
		}

		user, err := configuration.Queries.GetUserByEmail(r.Context(), requestBody.Email)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
		}

		printJsonResponse(w, mappers.MapUser(user, "", ""), http.StatusOK)
	}
}
