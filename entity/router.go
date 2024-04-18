package entity

import (
	"gorm.io/gorm"
)

type Router struct {
	gorm.Model
	Id         string `gorm:"<-:create;primary_key;type:varchar(36)"`
	Path       string `gorm:"type:varchar(64);not null;uniqueIndex:router_uni_idx"`
	Method     string `gorm:"type:varchar(8);not null;uniqueIndex:router_uni_idx"`
	Permission string `gorm:"type:varchar(16);not null"`
	// Foreign key
}
