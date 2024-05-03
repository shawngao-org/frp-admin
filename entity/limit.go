package entity

import (
	"gorm.io/gorm"
)

type Limit struct {
	gorm.Model
	Id       string `gorm:"<-:create;primary_key;type:varchar(36)"`
	UserId   string `gorm:"<-create;type:varchar(36);not null"`
	InBound  int64  `gorm:"type:int(5);not null"`
	OutBound int64  `gorm:"type:int(5);not null"`
	// Foreign key
}
