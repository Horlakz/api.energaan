package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/dto"
)

type CategoryDTO struct {
	dto.DTO
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	CreatedByID uuid.UUID `json:"createdBy"`
	UpdatedByID uuid.UUID `json:"updatedBy"`
	DeletedByID uuid.UUID `json:"deletedBy"`
}
