package service

import (
	"errors"
	"fmt"
	"frp-admin/db"
	"frp-admin/logger"
	"frp-admin/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func GetById(result interface{}, id any) error {
	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("result argument must be a pointer and must not be nil")
	}
	emptyEntity := reflect.New(val.Type()).Interface()
	db.Db.First(&result, "id = ?", id)
	if reflect.DeepEqual(result, emptyEntity) {
		return fmt.Errorf("%s not found", val.Type().Name())
	}
	return nil
}

func UpdateById(data interface{}) error {
	id, err := getIdFieldValue(data)
	if err != nil {
		return err
	}
	oldData := reflect.New(reflect.Indirect(reflect.ValueOf(data)).Type()).Interface()
	err = GetById(oldData, id.Interface())
	if err != nil {
		return err
	}
	util.CopyProperties(data, oldData)
	db.Db.Save(oldData)
	return nil
}

func DeleteById(e interface{}, ids ...any) error {
	val := reflect.ValueOf(e)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("result argument must be a pointer and must not be nil")
	}
	emptyEntity := reflect.New(val.Type()).Interface()
	db.Db.Delete(emptyEntity, ids)
	return nil
}

func getIdFieldValue(data interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return reflect.Value{}, errors.New("data argument must be a pointer and must not be nil")
	}
	val = val.Elem()
	idFieldName := "Id"
	idField := val.FieldByName(idFieldName)
	if !idField.CanInterface() {
		return reflect.Value{}, errors.New("id field cannot interface")
	}
	id := reflect.New(idField.Type()).Elem()
	id.Set(idField)
	return id, nil
}

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
