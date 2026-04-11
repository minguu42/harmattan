package atel

import (
	"log/slog"
	"reflect"
)

// maskAttr はlogタグの値がallowに設定されていない構造体フィールドの値をマスクする
// 非公開フィールドの値はゼロ値とする
// 本関数は以下の条件を満たす場合に実行される
// - 属性の値の種別が slog.KindAny である
// - 属性の値は構造体もしくは構造体のポインタである
func maskAttr(a slog.Attr) slog.Attr {
	if a.Value.Kind() != slog.KindAny {
		return a
	}

	switch rv := reflect.ValueOf(a.Value.Any()); rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
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
		newField := newStruct.Field(i)
		if !newField.CanSet() {
			continue
		}

		field := v.Field(i)
		if tag := structField.Tag.Get("log"); tag == "allow" {
			switch {
			case field.Kind() == reflect.Struct:
				newField.Set(maskStructValue(field))
			case field.Kind() == reflect.Pointer && !field.IsNil() && field.Elem().Kind() == reflect.Struct:
				newField.Set(maskStructValue(field.Elem()).Addr())
			default:
				newField.Set(field)
			}
			continue
		}
		setMaskedValue(newField)
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
	case reflect.Struct:
		structValue := reflect.New(t).Elem()
		for _, f := range structValue.Fields() {
			if f.CanSet() {
				setMaskedValue(f)
			}
		}
		newField.Set(structValue)
	default:
		newField.Set(reflect.Zero(t))
	}
}
