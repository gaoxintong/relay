package model

import (
	_ "gorm.io/driver/mysql" // gorm mysql 驱动包
	"gorm.io/gorm"
)

type user struct {
	gorm.Model
}
