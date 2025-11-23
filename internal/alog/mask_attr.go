package alog

import (
	"log/slog"
	"reflect"
)

// MaskAttr は以下の条件を満たす場合にlogタグの値がmaskに設定された構造体フィールドの値をマスクする
// - 属性の値の種別が slog.KindAny である
// - 属性の値は構造体もしくは構造体のポインタである
// - 構造体のすべてのフィールドがエクスポートされている
func MaskAttr(a slog.Attr) slog.Attr {
	if a.Value.Kind() != slog.KindAny {
		return a
	}

	switch rv := reflect.ValueOf(a.Value.Any()); rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() || rv.Elem().Kind() != reflect.Struct || !allFieldsExported(rv.Elem().Type()) {
			return a
		}
		return slog.Any(a.Key, maskStructValue(rv.Elem()).Addr().Interface())
	case reflect.Struct:
		if !allFieldsExported(rv.Type()) {
			return a
		}
		return slog.Any(a.Key, maskStructValue(rv).Interface())
	default:
		return a
	}
}

func allFieldsExported(t reflect.Type) bool {
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			return false
		}

		ft := f.Type
		switch ft.Kind() {
		case reflect.Pointer:
			if ft.Elem().Kind() == reflect.Struct && !allFieldsExported(ft.Elem()) {
				return false
			}
		case reflect.Struct:
			if !allFieldsExported(ft) {
				return false
			}
		default:
		}
	}
	return true
}

func maskStructValue(v reflect.Value) reflect.Value {
	t := v.Type()
	newStruct := reflect.New(t).Elem()
	for i := range v.NumField() {
		structField := t.Field(i)
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
