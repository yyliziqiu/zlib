package zutil

import (
	"fmt"
	"reflect"

	"github.com/yyliziqiu/zlib/zif"
)

func StructFields(model any) []string {
	mt := reflect.TypeOf(model)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		fields = append(fields, mt.Field(i).Name)
	}
	return fields
}

func StructValues(model any) []string {
	mv := reflect.ValueOf(model)
	var values []string
	for i := 0; i < mv.NumField(); i++ {
		values = append(values, fmt.Sprintf("%v", mv.Field(i).Interface()))
	}
	return values
}

func StructFieldValue(s any, fieldName string) (any, bool) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, false
	}
	return field.Interface(), true
}

func StructFieldTags(model any, tag string) []string {
	mt := reflect.TypeOf(model)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		tagval := f.Tag.Get(tag)
		fields = append(fields, zif.If(tagval != "", tagval, f.Name))
	}
	return fields
}
