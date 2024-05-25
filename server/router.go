package server

import (
	"fmt"
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
		user := v1.Group("/user")
		{
			user.POST("/login", api.Login)
			user.POST("/register", api.Register)
			user.POST("/verify-register", api.SendRegisterVerifyMail)
			user.POST("/confirm-register", api.ConfirmVerifyRegister)
			user.POST("/forget-password", api.SendForgetPasswordMail)
			user.POST("/reset-password", api.ResetPassword)
		}
		sys := v1.Group("/sys")
		{
			sys.POST("/ping", api.Ping)
		}
		docs.SwaggerInfo.Title = "API Docs"
		docs.SwaggerInfo.Description = "null."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%v", config.Conf.Server.Ip, config.Conf.Server.Port)
		docs.SwaggerInfo.BasePath = "/"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
