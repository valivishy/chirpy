package handlers

import (
	"chirpy/config"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/mappers"
	"chirpy/models"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func HandleLogin(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestBody := models.LoginUserRequest{}
		err := decoder.Decode(&requestBody)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := configuration.Queries.GetUserByEmail(r.Context(), requestBody.Email)
		if err != nil {
			respondWithError(w, "", http.StatusUnauthorized)
			return
		}

		if err := auth.CheckPasswordHash(user.HashedPassword, requestBody.Password); err != nil {
			respondWithError(w, "", http.StatusUnauthorized)
			return
		}

		jwt, err := auth.MakeJWT(user.ID, configuration.Secret)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = configuration.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		})
		if err != nil {
			respondWithError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		printJsonResponse(w, mappers.MapUser(user, jwt, refreshToken), http.StatusOK)
	}
}

func HandleRefresh(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromRefreshToken(configuration, r)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		jwt, err := auth.MakeJWT(user.ID, configuration.Secret)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}

		printJsonResponse(w, models.RefreshTokenResponse{Token: jwt}, http.StatusOK)
	}
}

func HandleRevoke(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if err = configuration.Queries.RevokeToken(r.Context(), token); err != nil {
			respondWithError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		printJsonResponse(w, nil, http.StatusNoContent)
	}
}

func getUserFromRefreshToken(configuration *config.Configuration, r *http.Request) (*database.User, error) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return nil, err
	}

	user, err := configuration.Queries.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		return nil, err
	}

	if len(user.Email) == 0 {
		return nil, errors.New("user not found")
	}

	return &user, err
}
