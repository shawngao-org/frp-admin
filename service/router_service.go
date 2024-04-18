package service

import (
	"errors"
	"frp-admin/db"
	"frp-admin/entity"
	"net/http"
	"reflect"
)

func CheckRouterPermission(req *http.Request, uid string) (bool, error) {
	path := req.URL.Path
	method := req.Method
	var router entity.Router
	db.Db.Table("router").Limit(1).Find(&router, "path = ? AND method = ?", path, method)
	if reflect.DeepEqual(router, entity.Router{}) {
		return false, errors.New("not found")
	}
	if router.Permission == "none" {
		return true, nil
	}
	if ExistRouterPermission(uid, router.Permission) {
		return true, nil
	}
	return false, errors.New("forbidden")
}
