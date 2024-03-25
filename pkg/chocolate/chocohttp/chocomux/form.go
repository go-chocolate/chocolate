package chocomux

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
)

func UnmarshalForm(form url.Values, data any) error {
	//if len(form) == 0 {
	//	return nil
	//}

	t := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("UnmarshalForm using unaddressable value")
	}

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		value = value.Elem()
	}

	switch t.Kind() {
	case reflect.Map:
		keyType := t.Key()
		if keyType.Kind() != reflect.String {
			return fmt.Errorf("the map key must be string when decode form to map, but got %v", t.Kind())
		}
		valType := t.Elem()
		for k, v := range form {
			val, err := set(valType, k, v)
			if err != nil {
				return err
			}
			value.SetMapIndex(reflect.ValueOf(k), val)
		}
		return nil
	case reflect.Struct:
		return unmarshalFormToStruct(&multipart.Form{Value: form, File: map[string][]*multipart.FileHeader{}}, t, value)
	default:
		return fmt.Errorf("cannot decode form as %v", t.Kind())
	}
}

func UnmarshalMultipartForm(form *multipart.Form, data any) error {
	//if len(form.File) == 0 && len(form.Value) == 0 {
	//	return nil
	//}

	t := reflect.TypeOf(data)
	value := reflect.ValueOf(data)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("UnmarshalForm using unaddressable value")
	}

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		value = value.Elem()
	}

	switch t.Kind() {
	case reflect.Map:
		keyType := t.Key()
		if keyType.Kind() != reflect.String {
			return fmt.Errorf("the map key must be string when decode form to map, but got %v", t.Kind())
		}
		valType := t.Elem()

		for k, v := range form.Value {
			val, err := set(valType, k, v)
			if err != nil {
				return err
			}
			value.SetMapIndex(reflect.ValueOf(k), val)
		}

		for k, v := range form.File {
			val, err := setFile(valType, k, v)
			if err != nil {
				return err
			}
			value.SetMapIndex(reflect.ValueOf(k), val)
		}

		return nil
	case reflect.Struct:
		return unmarshalFormToStruct(form, t, value)

	default:
		return fmt.Errorf("cannot decode form as %v", t.Kind())
	}
}

func unmarshalFormToStruct(form *multipart.Form, t reflect.Type, v reflect.Value) error {
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)
		for fieldValue.Kind() == reflect.Pointer {
			fieldValue = fieldValue.Elem()
		}
		if fieldType.Anonymous {
			if err := unmarshalFormToStruct(form, fieldType.Type, fieldValue); err != nil {
				return err
			}
		}
		name, must := structFieldName(fieldType, "form")
		_, ok1 := form.Value[name]
		_, ok2 := form.File[name]
		if must && !ok1 && !ok2 {
			return fmt.Errorf("missing form field '%s'", name)
		}
		var val reflect.Value
		var err error
		if ok1 {
			val, err = set(fieldType.Type, name, form.Value[name])
		} else if ok2 {
			val, err = setFile(fieldType.Type, name, form.File[name])
		} else {
			continue
		}
		if err != nil {
			return err
		}
		fieldValue.Set(val)
	}
	return nil
}

//func UnmarshalMultipartForm(form *multipart.Form, data any) error {
//	if len(form.File) == 0 && len(form.Value) == 0 {
//		return nil
//	}
//
//	t := reflect.TypeOf(data)
//	v := reflect.ValueOf(data)
//	if t.Kind() != reflect.Pointer {
//		return fmt.Errorf("UnmarshalMultipartForm using unaddressable value")
//	}
//
//	for t.Kind() == reflect.Pointer {
//		t = t.Elem()
//		v = v.Elem()
//	}
//	for i := 0; i < t.NumField(); i++ {
//		fieldType := t.Field(i)
//		fieldValue := v.Field(i)
//		name, must := structFieldName(fieldType, "form")
//
//		switch fieldValue.Kind() {
//		case reflect.Pointer:
//			if _, ok := form.File[name]; !ok {
//				if must {
//					return fmt.Errorf("missing form field '%s'", name)
//				} else {
//					continue
//				}
//			}
//			t1 := fieldType.Type.Elem()
//			if t1.Kind() != reflect.Struct || t1.PkgPath() != "mime/multipart" || t1.Name() != "FileHeader" {
//				return fmt.Errorf("cannot unmarshal field '%s'(%v) from multipart.Form", name, t1.Kind())
//			}
//
//			if files, ok := form.File[name]; !ok || len(files) == 0 {
//				return nil
//			} else {
//				fieldValue.Set(reflect.ValueOf(files[0]))
//			}
//		default:
//			if _, ok := form.Value[name]; !ok {
//				if must {
//					return fmt.Errorf("missing form field '%s'", name)
//				} else {
//					continue
//				}
//			}
//			val, err := set(fieldType.Type, name, form.Value[name])
//			if err != nil {
//				return err
//			}
//			fieldValue.Set(val)
//		}
//	}
//	return nil
//}

func set(t reflect.Type, key string, text []string) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Slice:
		return strings2slice(t, text)
	case reflect.Array:
		return strings2array(t, text)
	default:
		if len(text) == 0 {
			return reflect.Value{}, errors.New(fmt.Sprintf("no value for %s", key))
		}
		return string2value(t, text[0])
	}
}

func setFile(t reflect.Type, key string, fh []*multipart.FileHeader) (reflect.Value, error) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Slice:
		val := reflect.MakeSlice(t, 0, len(fh))
		for _, v := range fh {
			reflect.Append(val, reflect.ValueOf(v))
		}
		return val, nil
	case reflect.Array:
		val := reflect.New(t).Elem()
		for i := 0; i < val.Len() && i < len(fh); i++ {
			val.Index(i).Set(reflect.ValueOf(fh[i]))
		}
		return val, nil
	case reflect.Struct:
		if t.PkgPath() != "mime/multipart" || t.Name() != "FileHeader" {
			return reflect.Value{}, fmt.Errorf("cannot unmarshal field '%s'(%v) from multipart.Form", key, t.Kind())
		}
		fallthrough
	case reflect.Interface:
		return reflect.ValueOf(fh[0]), nil
	}
	return reflect.Value{}, fmt.Errorf("cannot unmarshal field '%s'(%v) from multipart.Form", key, t.Kind())
}

func toMultipartForm(f url.Values) *multipart.Form {
	form := &multipart.Form{
		Value: make(map[string][]string),
		File:  make(map[string][]*multipart.FileHeader),
	}
	return form
}

func strings2slice(t reflect.Type, valStrings []string) (reflect.Value, error) {
	val := reflect.MakeSlice(t, 0, len(valStrings))
	for _, v := range valStrings {
		item, err := string2value(t.Elem(), v)
		if err != nil {
			return reflect.Value{}, err
		}
		val = reflect.Append(val, item)
	}
	return val, nil
}

func strings2array(t reflect.Type, valStrings []string) (reflect.Value, error) {
	if t.Kind() != reflect.Array {
		return reflect.Value{}, errors.New("not array")
	}
	val := reflect.New(t).Elem()
	for i := 0; i < val.Len() && i < len(valStrings); i++ {
		data, err := string2value(t.Elem(), valStrings[i])
		if err != nil {
			return reflect.Value{}, err
		}
		val.Index(i).Set(data)
	}
	return val, nil
}

func string2value(t reflect.Type, valString string) (reflect.Value, error) {
	v, err := string2any(t, valString)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(v), nil
}

func string2any(t reflect.Type, valString string) (any, error) {
	kind := t.Kind()
	switch kind {
	case reflect.Bool:
	case reflect.Int:
		return toInt[int](valString)
	case reflect.Int8:
		return toInt[int8](valString)
	case reflect.Int16:
		return toInt[int16](valString)
	case reflect.Int32:
		return toInt[int32](valString)
	case reflect.Int64:
		return toInt[int64](valString)
	case reflect.Uint:
		return toInt[uint](valString)
	case reflect.Uint8:
		return toInt[uint8](valString)
	case reflect.Uint16:
		return toInt[uint16](valString)
	case reflect.Uint32:
		return toInt[uint32](valString)
	case reflect.Uint64:
		return toInt[uint64](valString)
	case reflect.Float32:
		return toFloat[float32](valString)
	case reflect.Float64:
		return toFloat[float64](valString)
	//case reflect.Array:
	//
	//case reflect.Slice:
	//
	case reflect.String:
		return valString, nil
	}
	return nil, fmt.Errorf("cannot convert %v to value", kind)
}

func toInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](text string) (T, error) {
	if i, err := strconv.ParseUint(text, 10, 64); err != nil {
		return 0, err
	} else {
		return T(i), nil
	}
}

func toIntVal[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](text string) (reflect.Value, error) {
	i, err := toInt[T](text)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(i), nil
}

func toFloat[T float32 | float64](text string) (T, error) {
	if f, err := strconv.ParseFloat(text, 64); err != nil {
		return 0, err
	} else {
		return T(f), nil
	}
}

func toFloatVal[T float32 | float64](text string) (reflect.Value, error) {
	i, err := toFloat[T](text)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(i), nil
}
