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

func HandleCreate(api *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestBody := models.CreateUserRequest{}
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

		user, err := api.Queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          requestBody.Email,
			HashedPassword: password,
		})
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
		}

		printJsonResponse(w, mappers.MapUser(user, "", ""), http.StatusCreated)
	}
}
