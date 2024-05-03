package service

import (
	"fmt"
	"frp-admin/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func OkHandleBySinglePayload(ctx *gin.Context, key string, value any) {
	ctx.JSON(http.StatusOK, gin.H{
		key: value,
	})
}

func OkHandleByPayload(ctx *gin.Context, payload map[string]any) {
	ctx.JSON(http.StatusOK, payload)
}

func OkHandle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
}

func ErrHandle(ctx *gin.Context, err any, status int) {
	if statusMsg := http.StatusText(status); statusMsg == "" {
		errMsg := fmt.Sprintf("%v is an invalid status code.", status)
		panic(errMsg)
	}
	if val, ok := err.(error); ok {
		ctx.JSON(status, gin.H{
			"message": val.Error(),
		})
		return
	}
	if val, ok := err.(string); ok {
		ctx.JSON(status, gin.H{
			"message": val,
		})
		return
	}
	errMsg := fmt.Sprintf("An argument of the wrong type appears to have been passed in. "+
		"At least the err argument should not be of type %s", reflect.TypeOf(err).Name())
	logger.LogErr(errMsg)
	panic(errMsg)
}
