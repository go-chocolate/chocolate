package types

import (
	"fmt"
	"strconv"
	"time"
)

func AnyToInt(val any) int {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, _ := strconv.Atoi(v)
		return i
	}
	return 0
}

func AnyToInt64(val any) int64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return int64(v)
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	}
	return 0
}

func AnyToInt32(val any) int32 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return int32(v)
	case int8:
		return int32(v)
	case int16:
		return int32(v)
	case int32:
		return int32(v)
	case int64:
		return int32(v)
	case uint:
		return int32(v)
	case uint8:
		return int32(v)
	case uint16:
		return int32(v)
	case uint32:
		return int32(v)
	case uint64:
		return int32(v)
	case float32:
		return int32(v)
	case float64:
		return int32(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return int32(i)
	}
	return 0
}

func AnyToInt16(val any) int16 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return int16(v)
	case int8:
		return int16(v)
	case int16:
		return int16(v)
	case int32:
		return int16(v)
	case int64:
		return int16(v)
	case uint:
		return int16(v)
	case uint8:
		return int16(v)
	case uint16:
		return int16(v)
	case uint32:
		return int16(v)
	case uint64:
		return int16(v)
	case float32:
		return int16(v)
	case float64:
		return int16(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return int16(i)
	}
	return 0
}

func AnyToInt8(val any) int8 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return int8(v)
	case int8:
		return int8(v)
	case int16:
		return int8(v)
	case int32:
		return int8(v)
	case int64:
		return int8(v)
	case uint:
		return int8(v)
	case uint8:
		return int8(v)
	case uint16:
		return int8(v)
	case uint32:
		return int8(v)
	case uint64:
		return int8(v)
	case float32:
		return int8(v)
	case float64:
		return int8(v)
	case string:
		i, _ := strconv.Atoi(v)
		return int8(i)
	}
	return 0
}

func AnyToUint(val any) uint {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return uint(v)
	case int8:
		return uint(v)
	case int16:
		return uint(v)
	case int32:
		return uint(v)
	case int64:
		return uint(v)
	case uint:
		return uint(v)
	case uint8:
		return uint(v)
	case uint16:
		return uint(v)
	case uint32:
		return uint(v)
	case uint64:
		return uint(v)
	case float32:
		return uint(v)
	case float64:
		return uint(v)
	case string:
		i, _ := strconv.ParseUint(v, 10, 64)
		return uint(i)
	}
	return 0
}

func AnyToString(val any) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func AnyToFloat64(val any) float64 {
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		i, _ := strconv.ParseFloat(v, 64)
		return i
	}
	return 0
}

func AnyToBool(val any) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	}
	return false
}

func AnyToTime(val any) time.Time {
	if val == nil {
		return time.Time{}
	}
	switch v := val.(type) {
	case time.Time:
		return v
	case string:
		t, _ := time.ParseInLocation(time.DateTime, v, time.Local)
		return t
	case fmt.Stringer:
		t, _ := time.ParseInLocation(time.DateTime, v.String(), time.Local)
		return t
	}
	return time.Time{}
}

func AnyTo[T any](a any) T {
	val, _ := a.(T)
	return val
}
