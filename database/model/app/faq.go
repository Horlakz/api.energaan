package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/model"
)

type Faq struct {
	model.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedByID uuid.UUID `gorm:"Column:created_by" json:"createdBy"`
	UpdatedByID uuid.UUID `gorm:"Column:updated_by" json:"updatedBy"`
	DeletedByID uuid.UUID `gorm:"Column:deleted_by" json:"deletedBy"`
}

func (Faq) TableName() string {
	return "faqs"
}
