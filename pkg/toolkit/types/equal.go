package types

import "reflect"

func deepEqual(a, b reflect.Value) bool {
	if a.Type().Name() != b.Type().Name() {
		return false
	}
	switch a.Kind() {
	case reflect.Pointer:
		if a.IsNil() && b.IsNil() {
			return true
		}
		if a.IsNil() || b.IsNil() {
			return false
		}
		return deepEqual(unptrValue(a), unptrValue(b))
	case reflect.Struct:
		for i := 0; i < a.NumField(); i++ {
			if !deepEqual(a.Field(i), b.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return false
		}
		for idx := 0; idx < a.Len(); idx++ {
			if !deepEqual(a.Index(idx), b.Index(idx)) {
				return false
			}
		}
		return true
	case reflect.Map:
		keys := a.MapKeys()
		for _, key := range keys {
			if !deepEqual(a.MapIndex(key), b.MapIndex(key)) {
				return false
			}
		}
		return true
	default:
		return a.Equal(b)
	}
}

// StructDeepEqual 比较两个结构体值是否相等
func StructDeepEqual(a, b any) bool {
	return deepEqual(reflect.ValueOf(a), reflect.ValueOf(b))
}
