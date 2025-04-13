package handlers

import (
	"chirpy/config"
	"chirpy/internal/database"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateUserRequest struct {
	Email string `json:"email"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func HandlerCreateUser(api *config.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestBody := CreateUserRequest{}
		err := decoder.Decode(&requestBody)
		if err != nil {
			respondWithError(w, "Something went wrong", http.StatusBadRequest)
			return
		}

		user, err := api.Queries.CreateUser(r.Context(), requestBody.Email)
		if err != nil {
			respondWithError(w, err.Error(), http.StatusBadRequest)
		}

		printJsonResponse(w, mapUser(user), http.StatusCreated)
	}
}

func mapUser(user database.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
