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

var Tables map[string]any

func Connect() {
	host := config.Conf.Database.Mysql.Host
	port := config.Conf.Database.Mysql.Port
	user := config.Conf.Database.Mysql.User
	pwd := config.Conf.Database.Mysql.Password
	db := config.Conf.Database.Mysql.Db
	logger.LogInfo("Connecting to mysql server [%s:%s]...", host, port)
	uri := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true", user, pwd, host, port, db)
	dsn, err := gorm.Open(mysql.Open(uri), &gorm.Config{})
	if err != nil {
		logger.LogErr("Can not connection to mysql server [%s:%s]", host, port)
		logger.LogErr("%s", err)
		os.Exit(-1)
	}
	logger.LogSuccess("Connected to mysql server.")
	connectMutex.Lock()
	Db = dsn
	Db = Db.Debug()
	connectMutex.Unlock()
	Tables = GetTableList()
	CheckTables()
}

func GetTableList() map[string]any {
	tables := make(map[string]any)
	tables["groups"] = &entity.Group{}
	tables["users"] = &entity.User{}
	tables["invites"] = &entity.Invite{}
	tables["nodes"] = &entity.Node{}
	tables["limits"] = &entity.Limit{}
	tables["proxies"] = &entity.Proxy{}
	tables["routers"] = &entity.Router{}
	tables["router_permissions"] = &entity.RouterPermission{}
	tables["tmp_codes"] = &entity.TmpCode{}
	return tables
}

func CheckTables() {
	for k, v := range Tables {
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

func ReinitializeDatabase() {
	for k, v := range Tables {
		if Db.Migrator().HasTable(k) {
			logger.LogWarn("Deleting table [%s] ...", k)
			err := Db.Migrator().DropTable(v)
			if err != nil {
				logger.LogErr("Delete table [%s] failed.", k)
				os.Exit(-1)
			}
			logger.LogSuccess("Table [%s] has been deleted.", k)
		}
	}
	CheckTables()
}
