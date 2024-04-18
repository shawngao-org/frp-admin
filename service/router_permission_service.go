package service

import (
	"frp-admin/db"
	"frp-admin/entity"
	"reflect"
)

func ExistRouterPermission(uid string, p string) bool {
	var rp entity.RouterPermission
	db.Db.Table("router_permission").Limit(1).Find(&rp, "user_id = ? AND permission = ?", uid, p)
	return !reflect.DeepEqual(rp, entity.RouterPermission{})
}
