package entity

import (
	"gorm.io/gorm"
	"time"
)

type Node struct {
	gorm.Model
	Id          string    `gorm:"<-:create;primary_key;type:varchar(36)"`
	Name        string    `gorm:"type:varchar(16);not null;uniqueIndex:name_idx"`
	HostName    string    `gorm:"type:varchar(32);not null"`
	Ip          string    `gorm:"type:varchar(39);not null"`
	Port        int64     `gorm:"type:int(5);not null"`
	AdminPort   int64     `gorm:"type:int(5);not null"`
	AdminPasswd string    `gorm:"type:varchar(64);not null"`
	Token       string    `gorm:"type:varchar(64);not null"`
	CreatedAt   time.Time `gorm:"<-:create;not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	Created     int64     `gorm:"<-:create;autoCreateTime;not null"`
	Updated     int64     `gorm:"autoUpdateTime:milli;not null"`
	// Foreign key
	Proxy Proxy `gorm:"foreignKey:NodeId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
