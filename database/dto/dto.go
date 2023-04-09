package dto

import (
	"time"

	"github.com/google/uuid"
)

type DTO struct {
	ID        uint      `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}
