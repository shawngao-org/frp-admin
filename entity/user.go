package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id           string    `gorm:"<-:create;primary_key;type:varchar(36)"`
	Name         string    `gorm:"type:varchar(16);not null;uniqueIndex:name_idx"`
	Email        string    `gorm:"type:varchar(32);not null;uniqueIndex:email_idx"`
	Password     string    `gorm:"type:varchar(128);not null"`
	TotpKey      string    `gorm:"type:varchar(64);not null"`
	IsValid      bool      `gorm:"default:false;not null"`
	RegisterTime time.Time `gorm:"<-:create;not null"`
	Ip           string    `gorm:"<-:create;type:varchar(39);not null;uniqueIndex:ip_idx"`
	Key          string    `gorm:"type:varchar(36);uniqueIndex:key_idx;not null"`
	GroupId      string    `gorm:"type:varchar(36);not null"`
	CreatedAt    time.Time `gorm:"<-:create;not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	Created      int64     `gorm:"<-:create;autoCreateTime;not null"`
	Updated      int64     `gorm:"autoUpdateTime:milli;not null"`
	// Foreign key
	Invite           []Invite         `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Invitees         Invite           `gorm:"foreignKey:Invitees;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Limit            Limit            `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Proxy            Proxy            `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	RouterPermission RouterPermission `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
