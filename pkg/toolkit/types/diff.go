package types

import (
	"reflect"
	"strings"
)

var JSONTagName = func(name string, tag reflect.StructTag) string {
	val := tag.Get("json")
	if n := strings.Index(val, ","); n > 0 {
		val = val[:n]
	}
	return val
}

// Diff 比较两个结构体字段差异
// 字段名默认取结构体字段名称，tagFunc 可自定义字段名
func Diff(a, b any, tagFunc ...func(name string, tag reflect.StructTag) string) map[string]any {
	t1 := unptrType(reflect.TypeOf(a))
	v1 := unptrValue(reflect.ValueOf(a))
	v2 := unptrValue(reflect.ValueOf(b))
	m := map[string]any{}
	for i := 0; i < v1.NumField(); i++ {
		vv1 := v1.Field(i).Interface()
		vv2 := v2.Field(i).Interface()
		if vv1 != vv2 {
			t := t1.Field(i)
			name := t.Name
			if len(tagFunc) > 0 {
				if val := tagFunc[0](name, t.Tag); val != "" {
					name = val
				}
			}
			m[name] = vv2
		}
	}
	return m
}
