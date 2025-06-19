package modifiers

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pchchv/modifier"
)

var (
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Duration(0))
)

// setValue allows setting of a specified value.
func setValueInner(field reflect.Value, param string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(param)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		if value, err := strconv.Atoi(param); err != nil {
			return err
		} else {
			field.SetInt(int64(value))
		}
	case reflect.Int64:
		var value int64
		if field.Type() == durationType {
			if d, err := time.ParseDuration(param); err != nil {
				return err
			} else {
				value = int64(d)
			}
		} else {
			if i, err := strconv.Atoi(param); err != nil {
				return err
			} else {
				value = int64(i)
			}
		}
		field.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value, err := strconv.Atoi(param); err != nil {
			return err
		} else {
			field.SetUint(uint64(value))
		}
	case reflect.Float32, reflect.Float64:
		if value, err := strconv.ParseFloat(param, 64); err != nil {
			return err
		} else {
			field.SetFloat(value)
		}
	case reflect.Bool:
		if value, err := strconv.ParseBool(param); err != nil {
			return err
		} else {
			field.SetBool(value)
		}
	case reflect.Map:
		var n int
		var err error
		if param != "" {
			if n, err = strconv.Atoi(param); err != nil {
				return err
			}
		}
		field.Set(reflect.MakeMapWithSize(field.Type(), n))
	case reflect.Slice:
		var cap int
		var err error
		if param != "" {
			if cap, err = strconv.Atoi(param); err != nil {
				return err
			}
		}
		field.Set(reflect.MakeSlice(field.Type(), 0, cap))
	case reflect.Struct:
		if field.Type() == timeType {
			if param != "" {
				if strings.ToLower(param) == "utc" {
					field.Set(reflect.ValueOf(time.Now().UTC()))
				} else {
					if t, err := time.Parse(time.RFC3339Nano, param); err != nil {
						return err
					} else {
						field.Set(reflect.ValueOf(t))
					}
				}
			} else {
				field.Set(reflect.ValueOf(time.Now()))
			}
		}
	case reflect.Chan:
		var buffer int
		var err error
		if param != "" {
			if buffer, err = strconv.Atoi(param); err != nil {
				return err
			}
		}
		field.Set(reflect.MakeChan(field.Type(), buffer))
	case reflect.Ptr:
		field.Set(reflect.New(field.Type().Elem()))
		return setValueInner(field.Elem(), param)
	}

	return nil
}

func setValue(_ context.Context, fl modifier.FieldLevel) error {
	return setValueInner(fl.Field(), fl.Param())
}

// defaultValue allows setting of a default value IF no value is already present.
func defaultValue(ctx context.Context, fl modifier.FieldLevel) error {
	if !fl.Field().IsZero() {
		return nil
	}
	return setValue(ctx, fl)
}
