package modifiers

import (
	"bytes"
	"context"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/pchchv/modifier"
	"github.com/segmentio/go-camelcase"
	"github.com/segmentio/go-snakecase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	nameRegex             = regexp.MustCompile(`[\p{L}]([\p{L}|[:space:]\-']*[\p{L}])*`)
	stripNumRegex         = regexp.MustCompile("[^0-9]")
	stripAlphaRegex       = regexp.MustCompile("[0-9]")
	stripAlphaUnicode     = regexp.MustCompile(`[\pL]`)
	stripNumUnicodeRegex  = regexp.MustCompile(`[^\pL]`)
	stripPunctuationRegex = regexp.MustCompile(`[[:punct:]]`)
	namePatterns          = []map[string]string{
		{`[^\pL-\s']`: ""}, // cut off everything except [ alpha, hyphen, whitespace, apostrophe]
		{`\s{2,}`: " "},    // trim more than two whitespaces to one
		{`-{2,}`: "-"},     // trim more than two hyphens to one
		{`'{2,}`: "'"},     // trim more than two apostrophes to one
		{`( )*-( )*`: "-"}, // trim enclosing whitespaces around hyphen
	}
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

// snakeCase converts string to snake case.
func snakeCase(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(snakecase.Snakecase(fl.Field().String()))
	}
	return nil
}

// slug converts string to a slug.
func slugCase(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(slug.Make(fl.Field().String()))
	}
	return nil
}

// titleCase converts string to title case,
// e.g. "this is a sentence" -> "This Is A Sentence".
func titleCase(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(cases.Title(language.Und, cases.NoLower).String(fl.Field().String()))
	}
	return nil
}

// stripAlphaCase removes all non-numeric characters.
// E. g.: "the price is €30,38" -> "3038".
// NOTE: The struct field will remain a string.
// No type conversion takes place.
func stripAlphaCase(_ context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(stripNumRegex.ReplaceAllLiteralString(fl.Field().String(), ""))
	}
	return nil
}

// stripNumCase removes all numbers.
// Example "39472349D34a34v69e8932747" -> "Dave".
// NOTE: The struct field will remain a string.
// No type conversion takes place.
func stripNumCase(_ context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(stripAlphaRegex.ReplaceAllLiteralString(fl.Field().String(), ""))
	}
	return nil
}

// stripNumUnicodeCase removes non-alpha unicode characters.
// E. g.: "!@£$%^&'()Hello 1234567890 World+[];\" -> "HelloWorld"
func stripNumUnicodeCase(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(stripNumUnicodeRegex.ReplaceAllLiteralString(fl.Field().String(), ""))
	}
	return nil
}

// stripAlphaUnicodeCase removes alpha unicode characters.
// E. g.: "Everything's here but the letters!" -> "' !".
func stripAlphaUnicodeCase(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(stripAlphaUnicode.ReplaceAllLiteralString(fl.Field().String(), ""))
	}
	return nil
}

// stripPunctuation removes punctuation.
// E. g.: "# M5W-1E6!!!" -> " M5W1E6".
func stripPunctuation(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(stripPunctuationRegex.ReplaceAllLiteralString(fl.Field().String(), ""))
	}
	return nil
}

// camelCase converts string to camel case.
func camelCase(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(camelcase.Camelcase(fl.Field().String()))
	}
	return nil
}

// nameCase Trims, strips numbers and special characters (except dashes and spaces separating names),
// converts multiple spaces and dashes to single characters, title cases multiple names.
// Example: "3493€848Jo-$%£@Ann " -> "Jo-Ann", " ~~ The Dude ~~" -> "The Dude", "**susan**" -> "Susan",
// " hugh fearnley-whittingstall" -> "Hugh Fearnley-Whittingstall".
func nameCase(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		fl.Field().SetString(cases.Title(language.Und, cases.NoLower).String(nameRegex.FindString(onlyOne(strings.ToLower(fl.Field().String())))))
	}
	return nil
}

func onlyOne(s string) string {
	for _, v := range namePatterns {
		for f, r := range v {
			s = regexp.MustCompile(f).ReplaceAllLiteralString(s, r)
		}
	}
	return s
}

func subStr(ctx context.Context, fl modifier.FieldLevel) error {
	switch fl.Field().Kind() {
	case reflect.String:
		val := fl.Field().String()
		params := strings.SplitN(fl.Param(), "-", 2)
		if len(params) == 0 || len(params[0]) == 0 {
			return nil
		}

		start, err := strconv.Atoi(params[0])
		if err != nil {
			return err
		}

		end := len(val)
		if len(params) >= 2 {
			if end, err = strconv.Atoi(params[1]); err != nil {
				return err
			}
		}

		if len(val) < start {
			fl.Field().SetString("")
			return nil
		}

		if len(val) < end {
			end = len(val)
		}

		if start > end {
			fl.Field().SetString("")
			return nil
		}

		fl.Field().SetString(val[start:end])
	}

	return nil
}
