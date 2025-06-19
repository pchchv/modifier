package modifiers

import (
	"context"
	"reflect"
	"strings"

	"github.com/pchchv/modifier"
)

// trimLeft trims extra left hand side of string using provided cutset.
func trimLeft(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(strings.TrimLeft(fl.Field().String(), fl.Param()))
	}
	return nil
}

// trimRight trims extra right hand side of string using provided cutset.
func trimRight(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(strings.TrimRight(fl.Field().String(), fl.Param()))
	}
	return nil
}
