package plan

import (
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database/model"
)

type Quote struct {
	model.Model
	FullName    string    `gorm:"Column:fullname" json:"fullname"`
	Email       string    `json:"email"`
	ServiceId   uuid.UUID `gorm:"Column:service_id" json:"serviceId"`
	ServiceType string    `gorm:"Column:service_type" json:"serviceType"`
	Phone       string    `json:"phone"`
	Country     string    `json:"country"`
}

func (Quote) TableName() string {
	return "quotes"
}
