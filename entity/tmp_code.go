package entity

import (
	"gorm.io/gorm"
)

type TmpCode struct {
	gorm.Model
	Id     string `gorm:"<-:create;primary_key;type:varchar(36)"`
	IsUsed bool   `gorm:"default:false;not null"`
	// Foreign key
}
