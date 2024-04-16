package entity

import (
	"gorm.io/gorm"
	"time"
)

type Group struct {
	gorm.Model
	Id            string    `gorm:"<-:create;primary_key;type:varchar(36)"`
	Name          string    `gorm:"type:varchar(16);not null;uniqueIndex:name_idx"`
	NickName      string    `gorm:"type:varchar(16);not null;uniqueIndex:nick_name_idx"`
	Traffic       int64     `gorm:"default:0;not null"`
	ProxyQuantity int64     `gorm:"default:0;not null"`
	BoundWidth    int64     `gorm:"default:0;not null"`
	CreatedAt     time.Time `gorm:"<-:create;not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	Created       int64     `gorm:"<-:create;autoCreateTime;not null"`
	Updated       int64     `gorm:"autoUpdateTime:milli;not null"`
	// Foreign key
	User []User `gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}
