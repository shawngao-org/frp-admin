package server

import (
	"frp-admin/api"
	"frp-admin/config"
	"frp-admin/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const apiPrefix = "/api/v1"

func LoadRouter(r *gin.Engine) {
	RootRouter(r)
}

func RootRouter(r *gin.Engine) {
	v1 := r.Group(apiPrefix)
	{
		v1.GET("/ping", api.Ping)
		v1.POST("/login", api.Login)
		v1.POST("/register", api.Register)

		docs.SwaggerInfo.Title = "API Docs"
		docs.SwaggerInfo.Description = "null."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = config.Conf.Server.Ip + ":" + config.Conf.Server.Port
		docs.SwaggerInfo.BasePath = "/"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
