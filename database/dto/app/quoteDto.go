package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/dto"
)

type QuoteDTO struct {
	dto.DTO
	FullName    string    `json:"fullname"`
	Email       string    `json:"email"`
	ServiceId   uuid.UUID `json:"serviceId"`
	ServiceType string    `json:"serviceType"`
	Phone       string    `json:"phone"`
	Country     string    `json:"country"`
}
