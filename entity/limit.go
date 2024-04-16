package entity

import (
	"gorm.io/gorm"
	"time"
)

type Limit struct {
	gorm.Model
	Id        string    `gorm:"<-:create;primary_key;type:varchar(36)"`
	UserId    string    `gorm:"<-create;type:varchar(36);not null"`
	InBound   int64     `gorm:"type:int(5);not null"`
	OutBound  int64     `gorm:"type:int(5);not null"`
	CreatedAt time.Time `gorm:"<-:create;not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Created   int64     `gorm:"<-:create;autoCreateTime;not null"`
	Updated   int64     `gorm:"autoUpdateTime:milli;not null"`
	// Foreign key
}
