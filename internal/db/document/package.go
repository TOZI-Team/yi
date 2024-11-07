package db_document

import (
	"github.com/ostafen/clover"
	"reflect"
	"strings"
)

// Package 定义了包的简略信息
type Package struct {
	Name        string // 包名称
	Author      string // 包作者
	Description string // 包的粗略解释
	Url         string // 包的详细表单位置
	Orphans     bool   // 是否已缺少支持
}

func ReflectMethod(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

func (receiver Package) ToMap() map[string]interface{} {
	t := reflect.TypeOf(receiver)
	v := reflect.ValueOf(receiver)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

// ToDocument 将一个 [Package] 转为 [clover.Document] 的指针
func (receiver Package) ToDocument() *clover.Document {
	d := clover.NewDocument()
	d.SetAll(receiver.ToMap())
	return d
}

func NewPackage() Package {
	return Package{}
}
