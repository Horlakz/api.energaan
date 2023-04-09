package auth

import (
	"github.com/horlakz/energaan-api/database/dto"
)

type UserDTO struct {
	dto.DTO
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
