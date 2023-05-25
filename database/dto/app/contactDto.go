package plan

import (
	"github.com/horlakz/energaan-api/database/dto"
)

type ContactDTO struct {
	dto.DTO
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
	Message  string `json:"message"`
}
