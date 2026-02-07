package env

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

func Load[T any]() (*T, error) {
	var v T
	rv := reflect.ValueOf(&v).Elem()
	if rv.Kind() != reflect.Struct {
		return nil, errors.New("type T must be a struct")
	}

	infos, err := gatherFieldInfos(rv)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		value, ok := os.LookupEnv(info.Key)
		if !ok {
			switch {
			case info.Default != "":
				value = info.Default
			case info.Required:
				return nil, fmt.Errorf("environment variable %s is required", info.Key)
			default:
				continue
			}
		}

		if err := processField(info.Field, value); err != nil {
			return nil, err
		}
	}
	return &v, nil
}

type fieldInfo struct {
	Field    reflect.Value
	Key      string
	Required bool
	Default  string
}

func gatherFieldInfos(v reflect.Value) ([]fieldInfo, error) {
	t := v.Type()

	infos := make([]fieldInfo, 0, v.NumField())
	for i := range v.NumField() {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}

		for f.Kind() == reflect.Pointer {
			if f.IsNil() {
				f.Set(reflect.New(f.Type().Elem()))
			}
			f = f.Elem()
		}

		fieldMeta := t.Field(i)
		info := fieldInfo{Field: f, Key: fieldMeta.Name}

		if tag := fieldMeta.Tag.Get("env"); tag != "" {
			if tag == "-" {
				continue
			}

			key, optionString, _ := strings.Cut(tag, ",")
			if key != "" {
				info.Key = key
			}

			options := strings.Split(optionString, ",")
			if slices.Contains(options, "required") {
				info.Required = true
			}
		}
		if tag := fieldMeta.Tag.Get("default"); tag != "" {
			info.Default = tag
		}
		infos = append(infos, info)

		if f.Kind() == reflect.Struct && textUnmarshaler(f) == nil {
			innerInfos, err := gatherFieldInfos(f)
			if err != nil {
				return nil, err
			}
			infos = append(infos[:len(infos)-1], innerInfos...)
		}
	}
	return infos, nil
}

func processField(field reflect.Value, value string) error {
	if unmarshaler := textUnmarshaler(field); unmarshaler != nil {
		return unmarshaler.UnmarshalText([]byte(value))
	}

	t := field.Type()
	switch field.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var v int64
		var err error
		if field.Kind() == reflect.Int64 && t.PkgPath() == "time" && t.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			v = int64(d)
		} else {
			v, err = strconv.ParseInt(value, 10, t.Bits())
		}
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, t.Bits())
		if err != nil {
			return err
		}
		field.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, t.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.String:
		field.SetString(value)
	case reflect.Slice:
		if t.Elem().Kind() != reflect.String {
			return fmt.Errorf("cannot handle slice of %s", t.Elem().Kind().String())
		}
		if value == "" {
			field.Set(reflect.MakeSlice(t, 0, 0))
		} else {
			parts := strings.Split(value, ",")
			slice := reflect.MakeSlice(t, len(parts), len(parts))
			for i, part := range parts {
				slice.Index(i).SetString(part)
			}
			field.Set(slice)
		}
	default:
		return fmt.Errorf("cannot handle field of kind %s", field.Kind().String())
	}
	return nil
}

func textUnmarshaler(v reflect.Value) encoding.TextUnmarshaler {
	if !v.CanInterface() {
		return nil
	}
	t, ok := v.Interface().(encoding.TextUnmarshaler)
	if !ok && v.CanAddr() {
		t, _ = v.Addr().Interface().(encoding.TextUnmarshaler)
	}
	return t
}
