package main

import (
	"frp-admin/config"
	"frp-admin/db"
	"frp-admin/redis"
)

func main() {
	config.GetConfig()
	db.Connect()
	redis.Connect()
}
