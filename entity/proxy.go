package entity

import (
	"gorm.io/gorm"
	"time"
)

type Proxy struct {
	gorm.Model
	Id             string    `gorm:"<-:create;primary_key;type:varchar(36)"`
	UserId         string    `gorm:"<-create;type:varchar(36);not null"`
	NodeId         string    `gorm:"<-create;type:varchar(36);not null"`
	ProxyName      string    `gorm:"type:varchar(16);not null;uniqueIndex:name_idx"`
	ProxyType      string    `gorm:"type:varchar(36);not null"`
	LocalIp        string    `gorm:"type:varchar(39);not null"`
	LocalPort      int64     `gorm:"type:int(5);not null"`
	UseEncryption  bool      `gorm:"default:false;not null"`
	UseCompression bool      `gorm:"default:false;not null"`
	Domain         string    `gorm:"type:varchar(36);uniqueIndex:domain_idx"`
	Router         string    `gorm:"type:varchar(128);"`
	RewriteHost    string    `gorm:"type:varchar(36)"`
	RemotePort     int64     `gorm:"type:int(5);not null;uniqueIndex:port_idx"`
	AccessKey      string    `gorm:"type:varchar(128)"`
	XFromWhere     string    `gorm:"type:varchar(36)"`
	CreatedAt      time.Time `gorm:"<-:create;not null"`
	UpdatedAt      time.Time `gorm:"not null"`
	Created        int64     `gorm:"<-:create;autoCreateTime;not null"`
	Updated        int64     `gorm:"autoUpdateTime:milli;not null"`
	// Foreign key
}
