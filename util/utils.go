package util

import "reflect"

func CopyProperties(source interface{}, target interface{}) {
	srcVal := reflect.ValueOf(source).Elem()
	dstVal := reflect.ValueOf(target).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		dstField := dstVal.Field(i)
		srcField := srcVal.Field(i)
		if dstField.CanSet() {
			dstField.Set(srcField)
		}
	}
}
