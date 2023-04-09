package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/model"
)

type Plan struct {
	model.Model
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	CreatedByID uuid.UUID `gorm:"Column:created_by" json:"createdBy"`
	UpdatedByID uuid.UUID `gorm:"Column:updated_by" json:"updatedBy"`
	DeletedByID uuid.UUID `gorm:"Column:deleted_by" json:"deletedBy"`
}

func (Plan) TableName() string {
	return "plans"
}
