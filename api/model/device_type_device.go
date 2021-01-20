package model

import (
	"gorm.io/gorm"
)

type DeviceTypeDevice struct {
	gorm.Model

	RecordID string `gorm:"unique:true"`

	DeviceID     string
	DeviceTypeID string
}
