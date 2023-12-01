package utils

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type DeepCopier interface {
	DeepCopy() interface{}
}

func DeepCopy[T any](v T) (c T, err error) {
	ptr := reflect.ValueOf(v)
	if ptr.IsNil() {
		return c, nil
	}
	if !isSinglePointer(ptr) {
		return c, errors.New("value is not a single pointer")
	}

	val := ptr.Elem()
	cp := reflect.New(val.Type()).Elem()

	if err := deepCopy(cp, val); err != nil {
		return c, fmt.Errorf("error copying value: %w", err)
	}

	return cp.Addr().Interface().(T), nil
}

func deepCopy(dst reflect.Value, src reflect.Value) error {
	if src.CanInterface() {
		if copier, ok := src.Interface().(DeepCopier); ok {
			dst.Set(reflect.ValueOf(copier.DeepCopy()))
			return nil
		}
	}

	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			return nil
		}
		val := src.Elem()

		if !val.IsValid() {
			return errors.New("invalid value")
		}
		dst.Set(reflect.New(val.Type()))
		if err := deepCopy(dst.Elem(), val); err != nil {
			return fmt.Errorf("error copying pointer: %w", err)
		}

	case reflect.Interface:
		if src.IsNil() {
			return nil
		}
		val := src.Elem()

		cpy := reflect.New(val.Type()).Elem()
		if err := deepCopy(cpy, val); err != nil {
			return fmt.Errorf("error copying interface: %w", err)
		}
		dst.Set(cpy)

	case reflect.Struct:
		t, ok := src.Interface().(time.Time)
		if ok {
			dst.Set(reflect.ValueOf(t))
			return nil
		}
		for i := 0; i < src.NumField(); i++ {
			// The Type's StructField for a given field is checked to see if StructField.PkgPath
			// is set to determine if the field is exported or not because CanSet() returns false
			// for settable fields.
			if src.Type().Field(i).PkgPath != "" {
				continue
			}
			if err := deepCopy(dst.Field(i), src.Field(i)); err != nil {
				return fmt.Errorf("error copying struct: %w", err)
			}
		}

	case reflect.Slice:
		if src.IsNil() {
			return nil
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			if err := deepCopy(dst.Index(i), src.Index(i)); err != nil {
				return fmt.Errorf("error copying slice: %w", err)
			}
		}

	case reflect.Map:
		if src.IsNil() {
			return nil
		}
		dst.Set(reflect.MakeMap(src.Type()))
		for _, key := range src.MapKeys() {
			val := src.MapIndex(key)
			copyValue := reflect.New(val.Type()).Elem()
			if err := deepCopy(copyValue, val); err != nil {
				return fmt.Errorf("error copying map value: %w", err)
			}
			copyKey := reflect.New(key.Type()).Elem()
			if err := deepCopy(copyKey, key); err != nil {
				return fmt.Errorf("error copying map key: %w", err)
			}
			dst.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}

	default:
		dst.Set(src)
	}
	return nil
}

func isSinglePointer(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.Elem().Kind() != reflect.Ptr
}
