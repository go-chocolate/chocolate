package jsonutil

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Type int8

const (
	Unknown Type = iota
	Null
	Text
	Number
	Boolean
	Object
	Array
)

var (
	null = []byte("null")
)

func (t Type) String() string {
	switch t {
	case Unknown:
		return "unknown"
	case Null:
		return "null"
	case Text:
		return "text"
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case Object:
		return "object"
	case Array:
		return "array"
	}
	return ""
}

var ErrUnknownType = errors.New("unknown field type")

type Any struct {
	value     string
	valueType Type
}

type AnyArray []*Any
type AnyMap map[string]*Any

func (a *Any) UnmarshalJSON(b []byte) error {
	val := string(b)
	if len(val) == 0 {
		return fmt.Errorf("unexpected json input")
	}
	if len(val) == 1 {
		if b[0] >= '0' && b[0] <= '9' {
			a.valueType = Number
			a.value = string(b)
			return nil
		} else {
			return fmt.Errorf("unexpected json input")
		}
	}
	if bytes.Equal(b, null) {
		a.valueType = Null
		a.value = string(b)
		return nil
	}

	if b[0] == '[' && b[len(b)-1] == ']' {
		a.valueType = Array
		a.value = string(b)
		return nil
	}

	if b[0] == '{' && b[len(b)-1] == '}' {
		a.valueType = Object
		a.value = string(b)
		return nil
	}

	if b[0] == '"' && b[len(b)-1] == '"' {
		a.valueType = Text
		a.value = string(b[1 : len(b)-1])
		return nil
	}
	if isBoolean(b) {
		a.valueType = Boolean
		a.value = string(b)
		return nil
	}
	if isNumber(b) {
		a.valueType = Number
		a.value = string(b)
		return nil
	}
	a.valueType = Unknown
	return fmt.Errorf("unexpected json input")
}

func isNumber(b []byte) bool {
	dot := false
	for _, v := range b {
		if v == '.' {
			if !dot {
				dot = true
				continue
			} else {
				return false
			}
		}
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func isBoolean(b []byte) bool {
	return bytes.Equal(b, []byte("true")) || bytes.Equal(b, []byte("false"))
}

func (a Any) MarshalJSON() ([]byte, error) {
	data := []byte(a.value)
	if a.valueType == Text {
		buf := make([]byte, len(a.value)+2)
		buf = append(buf, '"')
		buf = append(buf, data...)
		buf = append(buf, '"')
		return buf, nil
	}
	return data, nil
}

func (a *Any) Type() Type {
	return a.valueType
}

func (a *Any) Value() (any, error) {
	switch a.valueType {
	case Null:
		return nil, nil
	case Text:
		return a.value, nil
	case Number:
		if strings.Contains(a.value, ".") {
			v, _ := strconv.ParseFloat(a.value, 64)
			return v, nil
		} else {
			v, _ := strconv.ParseInt(a.value, 10, 64)
			return v, nil
		}
	case Boolean:
		return a.value == "true", nil
	case Object:
		var dst = make(map[string]*Any)
		if err := Unmarshal([]byte(a.value), &dst); err != nil {
			return nil, err
		} else {
			return dst, nil
		}
	case Array:
		var dst []*Any
		if err := Unmarshal([]byte(a.value), &dst); err != nil {
			return nil, err
		} else {
			return dst, nil
		}
	}
	return nil, ErrUnknownType
}

func (a *Any) Object() (map[string]*Any, error) {
	if a.valueType != Object {
		return nil, fmt.Errorf("content type is not object, but %v", a.valueType)
	}
	var dst = make(map[string]*Any)
	if err := Unmarshal([]byte(a.value), &dst); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}

func (a *Any) Array() ([]*Any, error) {
	if a.valueType != Array {
		return nil, fmt.Errorf("content type is not array, but %v", a.valueType)
	}
	var dst []*Any
	if err := Unmarshal([]byte(a.value), &dst); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}

func (a *Any) Int64() (int64, error) {
	if a.valueType != Number {
		return 0, fmt.Errorf("content type is not number, but %v", a.valueType)
	}
	if strings.Contains(a.value, ".") {
		val, err := a.Float64()
		return int64(val), err
	}
	return strconv.ParseInt(a.value, 10, 64)
}

func (a *Any) Float64() (float64, error) {
	if a.valueType != Number {
		return 0, fmt.Errorf("content type is not number, but %v", a.valueType)
	}
	return strconv.ParseFloat(a.value, 64)
}

func (a *Any) Boolean() (bool, error) {
	if a.valueType != Number {
		return false, fmt.Errorf("content type is not boolean, but %v", a.valueType)
	}
	return a.value == "true", nil
}

func (a *Any) String() string {
	if a.valueType == Null {
		return ""
	}
	return a.value
}

func (a *Any) IsNull() bool {
	return a.valueType == Null
}

func (a *Any) MustValue() any {
	val, err := a.Value()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustObject() map[string]*Any {
	val, err := a.Object()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustArray() []*Any {
	val, err := a.Array()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustInt64() int64 {
	val, err := a.Int64()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustFloat64() float64 {
	val, err := a.Float64()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustBoolean() bool {
	val, err := a.Boolean()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MaybeValue() any {
	val, err := a.Value()
	if err != nil {
		return nil
	}
	return val
}

func (a *Any) MaybeObject() map[string]*Any {
	val, err := a.Object()
	if err != nil {
		return nil
	}
	return val
}

func (a *Any) MaybeArray() []*Any {
	val, err := a.Array()
	if err != nil {
		return nil
	}
	return val
}

func (a *Any) MaybeInt64() int64 {
	val, err := a.Int64()
	if err != nil {
		return 0
	}
	return val
}

func (a *Any) MaybeFloat64() float64 {
	val, err := a.Float64()
	if err != nil {
		return 0
	}
	return val
}

func (a *Any) MaybeBoolean() bool {
	val, err := a.Boolean()
	if err != nil {
		return false
	}
	return val
}
