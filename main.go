package main

import (
	"frp-admin/config"
	"frp-admin/db"
	"frp-admin/logger"
	"frp-admin/redis"
	"frp-admin/server"
	"os"
)

type ExecutorType func(...any)

type Cmd struct {
	Description string
	Executor    ExecutorType
}

var Args map[string]Cmd

func main() {
	config.GetConfig()
	db.Connect()
	redis.Connect()
	Args = GetCliArgs()
	CommandHandler()
	server.HandleServer()
}

func GetCliArgs() map[string]Cmd {
	args := make(map[string]Cmd)
	args["-h"] = Cmd{"Show help info.", func(a ...any) {
		for k, v := range Args {
			logger.LogInfo("%-20s%-10s", k, v.Description)
		}
	}}
	args["--re-init-db"] = Cmd{"[Danger] Reinitialize the database.", func(a ...any) {
		db.ReinitializeDatabase()
	}}
	return args
}

func CommandHandler() {
	args := os.Args[1:]
	if len(args) <= 2 {
		return
	}
	for _, v := range args {
		if value, ok := Args[v]; ok {
			value.Executor()
		} else {
			logger.LogErr("Unknown command [%s].", v)
			os.Exit(-1)
		}
	}
	os.Exit(0)
}
