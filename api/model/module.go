package model

import (
	"gorm.io/gorm"
)

type Module struct {
	gorm.Model

	RecordID string `gorm:"unique:true"`

	Name   string
	Code   string `gorm:"unique:true"`
	AreaID string `gorm:"unique:true"`
}
