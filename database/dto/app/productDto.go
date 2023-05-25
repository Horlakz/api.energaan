package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/dto"
)

type ProductDTO struct {
	dto.DTO
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	CategoryID  uuid.UUID `json:"categoryId"`
	Images      []string  `json:"images"`
	Description string    `json:"description"`
	Features    []string  `json:"features"`
	CreatedByID uuid.UUID `json:"createdBy"`
	UpdatedByID uuid.UUID `json:"updatedBy"`
	DeletedByID uuid.UUID `json:"deletedBy"`
}
