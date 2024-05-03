package service

import (
	"errors"
	"frp-admin/db"
	"frp-admin/entity"
	"frp-admin/util"
	"reflect"
)

func GetTmpCodeById(id string) (entity.TmpCode, error) {
	var tmpCode entity.TmpCode
	db.Db.First(&tmpCode, "id = ?", id)
	if reflect.DeepEqual(tmpCode, entity.TmpCode{}) {
		return entity.TmpCode{}, errors.New("tmp code not found")
	}
	return tmpCode, nil
}

func UpdateTmpCodeById(tmpCode entity.TmpCode) error {
	result, err := GetTmpCodeById(tmpCode.Id)
	if err != nil {
		return err
	}
	util.CopyProperties(&tmpCode, &result)
	db.Db.Save(result)
	return nil
}
