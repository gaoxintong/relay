package model

import (
	"gorm.io/gorm"
)

type Area struct {
	gorm.Model

	RecordID string `gorm:"unique:true"`

	Name     string
	ParentID string
}
