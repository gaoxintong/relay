package model

import (
	"gorm.io/gorm"
)

type DeviceType struct {
	gorm.Model

	RecordID string `gorm:"unique:true"`

	Name string
}
