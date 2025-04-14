package models

import (
	"github.com/google/uuid"
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
