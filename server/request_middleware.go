package server

import (
	"frp-admin/config"
	"frp-admin/logger"
	"frp-admin/redis"
	"frp-admin/service"
	"frp-admin/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request pre logic code
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Conf.Develop {
			c.Next()
			return
		}
		token := c.Request.Header.Get("token")
		jt, err := util.ParseToken(token)
		if err != nil {
			logger.LogErr("%s", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message":   "Invalid token format.",
				"exception": err.Error(),
			})
			return
		}
		err = jt.Claims.Valid()
		if err != nil {
			logger.LogErr("%s", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message":   "Invalid token.",
				"exception": err.Error(),
			})
			return
		}
		uid := jt.Claims.(jwt.MapClaims)["user"].(string)
		exists := redis.Client.Exists(c, uid).Val()
		getToken := redis.Client.Get(c, uid).Val()
		if exists == 0 || getToken != token {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token.",
			})
			return
		}
		r, e := service.CheckRouterPermission(c.Request, uid)
		if e != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{
				"exception": e.Error(),
			})
			return
		}
		if !r {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		c.Next()
	}
}

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request after logic code
		// next
		c.Next()
	}
}
