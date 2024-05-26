package xj

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func SetDefaults(config any) error {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("config must be a pointer to a struct")
	}
	return setDefaults(v.Elem())
}

func setDefaults(v reflect.Value) error {
	for n := 0; n < v.NumField(); n++ {
		field := v.Field(n)
		if field.Kind() == reflect.Struct {
			if err := setDefaults(field); err != nil {
				return err
			}
		}
		tag := v.Type().Field(n).Tag.Get("default")
		if len(tag) > 0 && field.CanSet() {
			identification(tag, field)
		}
	}
	return nil
}

func identification(tag string, field reflect.Value) {
	switch field.Kind() {
	case reflect.Slice:
		slice := strings.Split(tag, ",")
		ns := reflect.MakeSlice(reflect.SliceOf(field.Type().Elem()), 0, len(slice))
		for _, elem := range slice {
			if v, err := identificationValue(field.Type().Elem().Kind(), elem); err == nil {
				ns = reflect.Append(ns, reflect.ValueOf(v).Convert(field.Type().Elem()))
			}
		}
		field.Set(ns)
	case reflect.Map:
		if field.IsNil() {
			field.Set(reflect.MakeMap(field.Type()))
		}
		for _, pair := range strings.Split(tag, ",") {
			kv := strings.Split(pair, ":")
			if len(kv) == 2 {
				if key, err := identificationValue(field.Type().Key().Kind(), kv[0]); err == nil {
					if value, err := identificationValue(field.Type().Elem().Kind(), kv[1]); err == nil {
						field.SetMapIndex(reflect.ValueOf(key).Convert(field.Type().Key()), reflect.ValueOf(value).Convert(field.Type().Elem()))
					}
				}
			}
		}
	case reflect.Interface:
		if val, err := identificationAny(tag); err == nil {
			field.Set(reflect.ValueOf(val))
		}
	default:
		if val, err := identificationValue(field.Kind(), tag); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	}
}

func identificationAny(value string) (any, error) {
	if v, err := strconv.Atoi(value); err == nil {
		return v, nil
	}
	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return v, nil
	}
	if v, err := strconv.ParseBool(value); err == nil {
		return v, nil
	}
	if strings.Contains(value, ":") {
		mv := make(map[any]any)
		for _, part := range strings.Split(value, ",") {
			if kv := strings.Split(part, ":"); len(kv) == 2 {
				if v, err := identificationAny(kv[1]); err == nil {
					mv[kv[0]] = v
				}
			}
		}
		return mv, nil
	}
	if strings.Contains(value, ",") {
		var slice []any
		for _, part := range strings.Split(value, ",") {
			if v, err := identificationAny(part); err == nil {
				slice = append(slice, v)
			}
		}
		return slice, nil
	}
	return value, nil
}

func identificationValue(kind reflect.Kind, value string) (any, error) {
	switch kind {
	case reflect.String:
		return value, nil
	case reflect.Int:
		return strconv.Atoi(value)
	case reflect.Int8:
		val, err := strconv.ParseInt(value, 10, 8)
		return int8(val), err
	case reflect.Int16:
		val, err := strconv.ParseInt(value, 10, 16)
		return int16(val), err
	case reflect.Int32:
		val, err := strconv.ParseInt(value, 10, 32)
		return int32(val), err
	case reflect.Int64:
		return strconv.ParseInt(value, 10, 64)
	case reflect.Uint:
		val, err := strconv.ParseUint(value, 10, 0)
		return uint(val), err
	case reflect.Uint8:
		val, err := strconv.ParseUint(value, 10, 8)
		return uint8(val), err
	case reflect.Uint16:
		val, err := strconv.ParseUint(value, 10, 16)
		return uint16(val), err
	case reflect.Uint32:
		val, err := strconv.ParseUint(value, 10, 32)
		return uint32(val), err
	case reflect.Uint64:
		return strconv.ParseUint(value, 10, 64)
	case reflect.Uintptr:
		val, err := strconv.ParseUint(value, 10, 0)
		return uintptr(val), err
	case reflect.Float32:
		val, err := strconv.ParseFloat(value, 32)
		return float32(val), err
	case reflect.Float64:
		return strconv.ParseFloat(value, 64)
	case reflect.Complex64:
		val, err := strconv.ParseComplex(value, 64)
		return complex64(val), err
	case reflect.Complex128:
		return strconv.ParseComplex(value, 128)
	case reflect.Bool:
		return strconv.ParseBool(value)
	}
	return nil, errors.New("undefined value type")
}
