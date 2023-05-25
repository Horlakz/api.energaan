package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/dto"
)

type FaqDTO struct {
	dto.DTO
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedByID uuid.UUID `json:"createdBy"`
	UpdatedByID uuid.UUID `json:"updatedBy"`
	DeletedByID uuid.UUID `json:"deletedBy"`
}
