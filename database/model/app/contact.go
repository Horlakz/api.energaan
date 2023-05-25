package plan

import (
	"github.com/horlakz/energaan-api/database/model"
)

type Contact struct {
	model.Model
	FullName string `gorm:"Column:fullname" json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
	Message  string `json:"message"`
}

func (Contact) TableName() string {
	return "contacts"
}
