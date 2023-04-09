package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/dto"
)

type PlanDTO struct {
	dto.DTO
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	CreatedByID uuid.UUID `json:"createdBy"`
	UpdatedByID uuid.UUID `json:"updatedBy"`
	DeletedByID uuid.UUID `json:"deletedBy"`
}
