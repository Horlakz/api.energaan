package auth

import (
	"github.com/horlakz/energaan-api/database/model"
)

type User struct {
	model.Model
	FullName string `json:"full_name"`
	Email    string `gorm:"UNIQUE_INDEX"`
	Password string `gorm:"Size:256"`
}

func (User) TableName() string {
	return "admins"
}
