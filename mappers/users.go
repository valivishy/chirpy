package mappers

import (
	"chirpy/internal/database"
	"chirpy/models"
)

func MapUser(user database.User, jwt string) models.UserDTO {
	return models.UserDTO{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwt,
	}
}
