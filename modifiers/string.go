package modifiers

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

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

// trimPrefix trims the string of a prefix.
func trimPrefix(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(strings.TrimPrefix(fl.Field().String(), fl.Param()))
	}
	return nil
}

// trimSuffix trims the string of a suffix.
func trimSuffix(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(strings.TrimSuffix(fl.Field().String(), fl.Param()))
	}
	return nil
}

// trimSpace trims extra space from text.
func trimSpace(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(strings.TrimSpace(fl.Field().String()))
	}
	return nil
}

// toLower convert string to lower case.
func toLower(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(strings.ToLower(fl.Field().String()))
	}
	return nil
}

// toUpper convert string to upper case.
func toUpper(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(strings.ToUpper(fl.Field().String()))
	}
	return nil
}

// uppercaseFirstCharacterCase converts a string so that it has only the first capital letter.
// E. g.: "all lower" -> "All lower".
func uppercaseFirstCharacterCase(_ context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		s := fl.Field().String()
		if s == "" {
			return nil
		}

		toRune, size := utf8.DecodeRuneInString(s)
		if !unicode.IsLower(toRune) {
			return nil
		}

		buf := &bytes.Buffer{}
		buf.WriteRune(unicode.ToUpper(toRune))
		buf.WriteString(s[size:])
		fl.Field().SetString(buf.String())
	}

	return nil
}
