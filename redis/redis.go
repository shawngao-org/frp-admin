package redis

import (
	"context"
	"fmt"
	"frp-admin/common"
	"frp-admin/config"
	"frp-admin/logger"
	"github.com/go-redis/redis/v8"
	"os"
	"sync"
	"time"
)

var (
	redisMutex sync.Mutex
	Client     *redis.Client
)

func Connect() {
	host := config.Conf.Database.Redis.Host
	port := config.Conf.Database.Redis.Port
	password := config.Conf.Database.Redis.Password
	db := config.Conf.Database.Redis.Db
	poolSize := config.Conf.Database.Redis.PoolSize
	timeout := config.Conf.Database.Redis.Timeout
	logger.LogInfo("Connecting to redis server [%s:%s]...", host, port)
	redisMutex.Lock()
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", host, port),
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})
	redisMutex.Unlock()
	ctx, cancel := context.WithTimeout(common.Context, time.Duration(timeout)*time.Second)
	defer cancel()
	res, err := Client.Ping(ctx).Result()
	if err != nil {
		logger.LogErr("Can not connection to redis server [%s:%s]", host, port)
		logger.LogErr("%s", err)
		os.Exit(-1)
	}
	logger.LogSuccess("Ping => %v", res)
	logger.LogSuccess("Connected to redis server.")
}
