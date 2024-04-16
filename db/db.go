package db

import (
	"fmt"
	"frp-admin/config"
	"frp-admin/entity"
	"frp-admin/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"sync"
)

var (
	connectMutex sync.Mutex
	Db           *gorm.DB
)

func Connect() {
	host := config.Conf.Database.Mysql.Host
	port := config.Conf.Database.Mysql.Port
	user := config.Conf.Database.Mysql.User
	pwd := config.Conf.Database.Mysql.Password
	db := config.Conf.Database.Mysql.Db
	logger.LogInfo("Connecting to mysql server [%s:%s]...", host, port)
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, db)
	dsn, err := gorm.Open(mysql.Open(uri), &gorm.Config{})
	if err != nil {
		logger.LogErr("Can not connection to mysql server [%s:%s]", host, port)
		logger.LogErr("%s", err)
		os.Exit(-1)
	}
	logger.LogSuccess("Connected to mysql server.")
	connectMutex.Lock()
	Db = dsn
	connectMutex.Unlock()
	tables := GetTableList()
	CheckTables(tables)
}

func GetTableList() map[string]any {
	tables := make(map[string]any)
	tables["group"] = &entity.Group{}
	tables["users"] = &entity.User{}
	tables["invites"] = &entity.Invite{}
	tables["nodes"] = &entity.Node{}
	tables["limits"] = &entity.Limit{}
	tables["proxies"] = &entity.Proxy{}
	return tables
}

func CheckTables(tables map[string]any) {
	for k, v := range tables {
		if !Db.Migrator().HasTable(k) {
			logger.LogWarn("Table [%s] does not exists. Creating...", k)
			err := Db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(v)
			if err != nil {
				logger.LogErr("Create table [%s] failed.", k)
				os.Exit(-1)
			}
			logger.LogSuccess("Table [%s] has been created.", k)
		}
	}
}
