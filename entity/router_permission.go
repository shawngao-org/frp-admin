package entity

import (
	"gorm.io/gorm"
)

type RouterPermission struct {
	gorm.Model
	Id         string `gorm:"<-:create;primary_key;type:varchar(36)"`
	UserId     string `gorm:"<-create;type:varchar(36);not null"`
	Permission string `gorm:"type:varchar(16);not null"`
	// Foreign key
}
