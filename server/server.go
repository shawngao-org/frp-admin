package server

import (
	"context"
	"frp-admin/config"
	"frp-admin/logger"
	"frp-admin/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func HandleServer() {
	r := gin.Default()
	r.Use(RequestMiddleware(), CORSMiddleware(), ResponseMiddleware())
	LoadRouter(r)
	for _, router := range r.Routes() {
		logger.LogInfo("[%s] => [%s]", router.Method, router.Path)
	}
	srv := &http.Server{
		Addr:    config.Conf.Server.Ip + ":" + config.Conf.Server.Port,
		Handler: r,
	}
	go func() {
		startServer(srv)
	}()
	shutdownServer(srv)
}

func startServer(srv *http.Server) {
	logger.LogInfo("Starting server...")
	logger.LogSuccess("Listening address: [%s]", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		logger.LogErr("Error: %s", err)
		os.Exit(-1)
	}
}

func shutdownServer(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.LogWarn("Shutting down server...")
	err := redis.Client.Close()
	if err != nil {
		logger.LogErr("Error: %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		logger.LogErr("Error: %s", err)
	}
	logger.LogSuccess("The server has been closed.")
}
