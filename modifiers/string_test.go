package modifiers

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var excludedStructs = map[reflect.Type]bool{
	reflect.TypeOf(time.Time{}):      true,
	reflect.TypeOf(sql.NullString{}): true,
}

func TestMultiple(t *testing.T) {
	assert := require.New(t)
	conform := New()
	s := interface{}("PCHCHV ")
	err := conform.Field(context.Background(), &s, "trim,lcase")
	assert.NoError(err)
	assert.Equal("pchchv", s)
}

func TestEnumType(t *testing.T) {
	assert := require.New(t)

	type State string
	const START State = "start"

	state := State("START")
	conform := New()
	err := conform.Field(context.Background(), &state, "lcase")
	assert.NoError(err)
	assert.Equal(START, state)
}

func TestEmails(t *testing.T) {
	conform := New()
	email := "           Jack.Pochechuev@gmail.com            "

	type Test struct {
		Email string `mod:"trim"`
	}

	tt := Test{Email: email}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}

	if tt.Email != "Jack.Pochechuev@gmail.com" {
		t.Fatalf("Unexpected value '%s'\n", tt.Email)
	}

	if err := conform.Field(context.Background(), &email, "trim"); err != nil {
		log.Fatal(err)
	}
	if email != "Jack.Pochechuev@gmail.com" {
		t.Fatalf("Unexpected value '%s'\n", tt.Email)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "trim"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = "    test     "
	if err := conform.Field(context.Background(), &iface, "trim"); err != nil {
		log.Fatal(err)
	}
	if iface != "test" {
		t.Fatalf("Unexpected value '%s'\n", "test")
	}
}

func TestTrimLeft(t *testing.T) {
	conform := New()
	s := "#$%_test"
	expected := "test"

	type Test struct {
		String string `mod:"ltrim=#_$%"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "ltrim=%_$#"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "ltrim=%_$#"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "ltrim=%_$#"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestTrimRight(t *testing.T) {
	conform := New()
	s := "test#$%_"
	expected := "test"

	type Test struct {
		String string `mod:"rtrim=#_$%"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "rtrim=#_$%"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "rtrim=#_$%"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "rtrim=#_$%"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestTrimPrefix(t *testing.T) {
	conform := New()
	s := "pre-test"
	expected := "test"

	type Test struct {
		String string `mod:"tprefix=pre-"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "tprefix=pre-"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "tprefix=pre-"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "tprefix=pre-"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestTrimSuffix(t *testing.T) {
	conform := New()
	s := "test-suffix"
	expected := "test"

	type Test struct {
		String string `mod:"tsuffix=-suffix"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "tsuffix=-suffix"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "tsuffix=-suffix"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "tsuffix=-suffix"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestToLower(t *testing.T) {
	conform := New()
	s := "TEST"
	expected := "test"

	type Test struct {
		String string `mod:"lcase"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "lcase"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "lcase"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "lcase"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestToUpper(t *testing.T) {
	conform := New()
	s := "test"
	expected := "TEST"

	type Test struct {
		String string `mod:"ucase"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "ucase"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "ucase"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "ucase"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestSnakeCase(t *testing.T) {
	conform := New()
	s := "ThisIsSNAKEcase"
	expected := "this_is_snakecase"

	type Test struct {
		String string `mod:"snake"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "snake"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "snake"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "snake"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestTitleCase(t *testing.T) {
	conform := New()
	s := "this is a sentence"
	expected := "This Is A Sentence"

	type Test struct {
		String string `mod:"title"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "title"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "title"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "title"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestSlugCase(t *testing.T) {
	conform := New()
	s := "this-is +a SentencE9"
	expected := "this-is-a-sentence9"

	type Test struct {
		String string `mod:"slug"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "slug"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "slug"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "slug"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestNameCase(t *testing.T) {
	conform := New()
	s := "3493€848Jo-$%£@Ann "
	expected := "Jo-Ann"

	type Test struct {
		String string `mod:"name"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "name"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "name"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "name"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}

	s = " ~~ The Dude ~~"
	expected = "The Dude"
	if err := conform.Field(context.Background(), &s, "name"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	s = "**susan**"
	expected = "Susan"
	if err := conform.Field(context.Background(), &s, "name"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	s = " hugh fearnley-whittingstall"
	expected = "Hugh Fearnley-Whittingstall"
	if err := conform.Field(context.Background(), &s, "name"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}
}

func TestUCFirstCase(t *testing.T) {
	conform := New()
	s := "this is uc first case"
	expected := "This is uc first case"

	type Test struct {
		String string `mod:"ucfirst"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "ucfirst"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "ucfirst"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "ucfirst"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}

	s = ""
	expected = ""
	if err := conform.Field(context.Background(), &s, "ucfirst"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}
}

func TestNumCase(t *testing.T) {
	conform := New()
	s := "the price is €30,38"
	expected := "3038"

	type Test struct {
		String string `mod:"strip_alpha"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "strip_alpha"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "strip_alpha"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s

	if err := conform.Field(context.Background(), &iface, "strip_alpha"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestNotNumCase(t *testing.T) {
	conform := New()
	s := "39472349D34a34v69e8932747"
	expected := "Dave"

	type Test struct {
		String string `mod:"strip_num"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "strip_num"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "strip_num"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "strip_num"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestAlphaCase(t *testing.T) {
	conform := New()
	s := "!@£$%^&'()Hello 1234567890 World+[];\\"
	expected := "HelloWorld"

	type Test struct {
		String string `mod:"strip_num_unicode"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "strip_num_unicode"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "strip_num_unicode"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "strip_num_unicode"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestNotAlphaCase(t *testing.T) {
	conform := New()
	s := "Everything's here but the letters!"
	expected := "'    !"

	type Test struct {
		String string `mod:"strip_alpha_unicode"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "strip_alpha_unicode"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "strip_alpha_unicode"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "strip_alpha_unicode"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestPunctuation(t *testing.T) {
	conform := New()
	s := "# M5W-1E6!!!"
	expected := " M5W1E6"

	type Test struct {
		String string `mod:"strip_punctuation"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "strip_punctuation"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "strip_punctuation"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "strip_punctuation"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestCamelCase(t *testing.T) {
	conform := New()
	s := "this_is_snakecase"
	expected := "thisIsSnakecase"

	type Test struct {
		String string `mod:"camel"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	if err := conform.Field(context.Background(), &s, "camel"); err != nil {
		log.Fatal(err)
	}
	if s != expected {
		t.Fatalf("Unexpected value '%s'\n", s)
	}

	var iface interface{}
	if err := conform.Field(context.Background(), &iface, "camel"); err != nil {
		log.Fatal(err)
	}
	if iface != nil {
		t.Fatalf("Unexpected value '%v'\n", nil)
	}

	iface = s
	if err := conform.Field(context.Background(), &iface, "camel"); err != nil {
		log.Fatal(err)
	}
	if iface != expected {
		t.Fatalf("Unexpected value '%v'\n", iface)
	}
}

func TestString(t *testing.T) {
	assert := require.New(t)
	conform := New()
	conform.RegisterInterceptor(func(current reflect.Value) (inner reflect.Value) {
		current.FieldByName("Valid").SetBool(true)
		return current.FieldByName("String")
	}, sql.NullString{})

	tests := []struct {
		name        string
		field       interface{}
		tags        string
		expected    interface{}
		expectError bool
	}{
		{
			name:     "test Camel Case",
			field:    "this_is_snakecase",
			tags:     "camel",
			expected: "thisIsSnakecase",
		},
		{
			name: "test Camel Case struct",
			field: struct {
				String string `mod:"camel"`
			}{String: "this_is_snakecase"},
			tags: "camel",
			expected: struct {
				String string `mod:"camel"`
			}{String: "thisIsSnakecase"},
		},
		{
			name:     "sql.nullString lcase",
			field:    sql.NullString{String: "UPPERCASE", Valid: true},
			tags:     "lcase",
			expected: sql.NullString{String: "uppercase", Valid: true},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var err error

			input := reflect.ValueOf(tc.field)
			if !input.CanAddr() {
				// create NEW addressable pointer to struct and assign the old unadressable one.
				// sort fo like:
				//
				// var newStruct *oldstructdef
				// *newStruct = *&oldStruct
				//
				newVal := reflect.New(input.Type())
				newVal.Elem().Set(input)
				tc.field = newVal.Interface()
			}

			if reflect.ValueOf(tc.field).Kind() == reflect.Struct && !excludedStructs[reflect.TypeOf(tc.field)] {
				err = conform.Struct(context.Background(), &tc.field)
			} else {
				err = conform.Field(context.Background(), &tc.field, tc.tags)
			}

			if tc.expectError {
				assert.Error(err)
				return
			}

			assert.NoError(err)
			assert.Equal(tc.expected, reflect.Indirect(reflect.ValueOf(tc.field)).Interface())
		})
	}
}

func TestSubStr(t *testing.T) {
	conform := New()
	s := "123"
	expected := "123"

	type Test struct {
		String string `mod:"substr=0-3"`
	}

	tt := Test{String: s}
	if err := conform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	if tt.String != expected {
		t.Fatalf("Unexpected value '%s'\n", tt.String)
	}

	tag := "substr=f-3"
	if err := conform.Field(context.Background(), &s, tag); err == nil {
		t.Fatalf("Unexpected value '%s' instead of error for tag %s\n", s, tag)
	}
	tag = "substr=2-f"
	if err := conform.Field(context.Background(), &s, tag); err == nil {
		t.Fatalf("Unexpected value '%s' instead of error for tag %s\n", s, tag)
	}

	tests := []struct {
		tag      string
		expected string
	}{
		{
			tag:      "substr",
			expected: "123",
		},
		{
			tag:      "substr=0-1",
			expected: "1",
		},
		{
			tag:      "substr=0-3",
			expected: "123",
		},
		{
			tag:      "substr=0-2",
			expected: "12",
		},
		{
			tag:      "substr=1-2",
			expected: "2",
		},
		{
			tag:      "substr=3-3",
			expected: "",
		},
		{
			tag:      "substr=4-5",
			expected: "",
		},
		{
			tag:      "substr=2-1",
			expected: "",
		},
		{
			tag:      "substr=2-5",
			expected: "3",
		},
		{
			tag:      "substr=2",
			expected: "3",
		},
	}
	for _, test := range tests {
		st := s

		if err := conform.Field(context.Background(), &st, test.tag); err != nil {
			log.Fatal(err)
		}
		if st != test.expected {
			t.Fatalf("Unexpected value '%s' for tag %s\n", st, test.tag)
		}

		var iface interface{}
		if err := conform.Field(context.Background(), &iface, test.tag); err != nil {
			log.Fatal(err)
		}
		if iface != nil {
			t.Fatalf("Unexpected value '%v'\n", nil)
		}

		iface = s
		if err := conform.Field(context.Background(), &iface, test.tag); err != nil {
			log.Fatal(err)
		}
		if iface != test.expected {
			t.Fatalf("Unexpected value '%v'\n", iface)
		}
	}
}
