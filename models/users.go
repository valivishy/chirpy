package models

import (
	"github.com/google/uuid"
	"time"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDTO struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type LoginUserRequest struct {
	Email            string  `json:"email"`
	Password         string  `json:"password"`
	ExpiresInSeconds *string `json:"expires_in_seconds"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
