package chocomux

import (
	"reflect"
	"strings"
)

func structFieldName(t reflect.StructField, tagName string) (string, bool) {
	must := t.Tag.Get("required") == "true"
	name := t.Name
	tag := t.Tag.Get(tagName)
	if tag == "" {
		return name, must
	}
	temp := strings.TrimSpace(strings.Split(tag, ",")[0])
	if temp != "" {
		return temp, must
	}
	return name, must
}
