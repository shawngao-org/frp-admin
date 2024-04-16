package entity

import (
	"gorm.io/gorm"
	"time"
)

type Invite struct {
	gorm.Model
	Id        string    `gorm:"<-:create;primary_key;type:varchar(36)"`
	UserId    string    `gorm:"<-create;type:varchar(36);not null"`
	code      string    `gorm:"<-create;type:varchar(36);uniqueIndex:code_idx"`
	Invitees  string    `gorm:"<-create;type:varchar(36);uniqueIndex:invite_idx"`
	CreatedAt time.Time `gorm:"<-:create;not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Created   int64     `gorm:"<-:create;autoCreateTime;not null"`
	Updated   int64     `gorm:"autoUpdateTime:milli;not null"`
	// Foreign key
}
