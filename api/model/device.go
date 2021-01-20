package model

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model

	RecordID string `gorm:"unique:true"`

	Name     string
	ModuleID string
	RouteNO  int
	InputNO  int
}
