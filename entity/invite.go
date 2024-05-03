package entity

import (
	"gorm.io/gorm"
)

type Invite struct {
	gorm.Model
	Id       string `gorm:"<-:create;primary_key;type:varchar(36)"`
	UserId   string `gorm:"<-create;type:varchar(36);not null"`
	Code     string `gorm:"<-create;type:varchar(36);uniqueIndex:code_idx"`
	Invitees string `gorm:"<-create;type:varchar(36);uniqueIndex:invite_idx"`
	// Foreign key
}
