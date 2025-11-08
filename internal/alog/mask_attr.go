package alog

import (
	"log/slog"
	"reflect"
)

// MaskAttr masks sensitive fields in log attributes
func MaskAttr(a slog.Attr) slog.Attr {
	if a.Value.Kind() != slog.KindAny {
		return a
	}

	switch rv := reflect.ValueOf(a.Value.Any()); rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
			return a
		}
		// TODO(furukawa): errtrace.StackError でも正しく動くように修正する
		if rv.Type().Elem().Name() == "StackError" {
			return a
		}
		return slog.Any(a.Key, maskStructValue(rv.Elem()).Addr().Interface())
	case reflect.Struct:
		return slog.Any(a.Key, maskStructValue(rv).Interface())
	default:
		return a
	}
}

func maskStructValue(v reflect.Value) reflect.Value {
	t := v.Type()
	newStruct := reflect.New(t).Elem()
	for i := range v.NumField() {
		structField := t.Field(i)
		if !structField.IsExported() {
			continue
		}
		newField := newStruct.Field(i)
		if !newField.CanSet() {
			continue
		}

		if tag := structField.Tag.Get("log"); tag == "mask" {
			setMaskedValue(newField)
			continue
		}

		field := v.Field(i)
		switch {
		case field.Kind() == reflect.Struct:
			newField.Set(maskStructValue(field))
		case field.Kind() == reflect.Pointer && !field.IsNil() && field.Elem().Kind() == reflect.Struct:
			newField.Set(maskStructValue(field.Elem()).Addr())
		default:
			newField.Set(field)
		}
	}
	return newStruct
}

func setMaskedValue(newField reflect.Value) {
	t := newField.Type()
	switch t.Kind() {
	case reflect.Bool:
		newField.SetBool(false)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		newField.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		newField.SetUint(0)
	case reflect.Float32, reflect.Float64:
		newField.SetFloat(0.0)
	case reflect.Map:
		newField.Set(reflect.MakeMap(t))
	case reflect.Slice:
		newField.Set(reflect.MakeSlice(t, 0, 0))
	case reflect.String:
		newField.SetString("<hidden>")
	default:
		newField.Set(reflect.Zero(t))
	}
}
