package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	RecordID string `gorm:"unique:true"`

	Name     string
	PassCode string
	NickName string
	Role     string
}
