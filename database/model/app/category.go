package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/model"
)

type Category struct {
	model.Model
	Slug        string    `gorm:"UNIQUE_INDEX" json:"slug"`
	Name        string    `json:"name"`
	CreatedByID uuid.UUID `gorm:"Column:created_by" json:"createdBy"`
	UpdatedByID uuid.UUID `gorm:"Column:updated_by" json:"updatedBy"`
	DeletedByID uuid.UUID `gorm:"Column:deleted_by" json:"deletedBy"`
}

func (Category) TableName() string {
	return "categories"
}
