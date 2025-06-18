package modifier

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	. "github.com/pchchv/go-assert"
)

func TestBadValues(t *testing.T) {
	tform := New()
	tform.Register("blah", func(ctx context.Context, fl FieldLevel) error { return nil })

	type Test struct {
		Ignore string `mold:"-"`
		String string `mold:"blah,,blah"`
	}

	var tt Test
	err := tform.Struct(context.Background(), &tt)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "invalid tag '' found on field String")

	err = tform.Struct(context.Background(), tt)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "mold: Struct(non-pointer mold.Test)")

	err = tform.Struct(context.Background(), nil)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "mold: Struct(nil)")

	var i int
	err = tform.Struct(context.Background(), &i)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "mold: (nil *int)")

	var iface interface{}
	err = tform.Struct(context.Background(), iface)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "mold: Struct(nil)")

	iface = nil
	err = tform.Struct(context.Background(), &iface)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "mold: (nil *interface {})")

	var tst *Test
	err = tform.Struct(context.Background(), tst)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "mold: Struct(nil *mold.Test)")

	var tm *time.Time
	err = tform.Field(context.Background(), tm, "blah")
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "mold: Field(nil *time.Time)")

	PanicMatches(t, func() { tform.Register("", nil) }, "Function Key cannot be empty")
	PanicMatches(t, func() { tform.Register("test", nil) }, "Function cannot be empty")
	PanicMatches(t, func() {
		tform.Register(",", func(ctx context.Context, fl FieldLevel) error { return nil })
	}, "Tag ',' either contains restricted characters or is the same as a restricted tag needed for normal operation")

	PanicMatches(t, func() { tform.RegisterAlias("", "") }, "Alias cannot be empty")
	PanicMatches(t, func() { tform.RegisterAlias("test", "") }, "Aliased tags cannot be empty")
	PanicMatches(t, func() { tform.RegisterAlias(",", "test") }, "Alias ',' either contains restricted characters or is the same as a restricted tag needed for normal operation")
}

func TestDiveKeys(t *testing.T) {
	type Test struct {
		Map map[string]string `s:"dive,keys,default,endkeys,default"`
	}

	set := New()
	set.SetTagName("s")
	set.Register("default", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().Set(reflect.ValueOf("after"))
		return nil
	})
	set.Register("err", func(ctx context.Context, fl FieldLevel) error {
		return errors.New("err")
	})

	test := Test{
		Map: map[string]string{
			"b4": "b4",
		},
	}

	err := set.Struct(context.Background(), &test)
	Equal(t, err, nil)

	val := test.Map["after"]
	Equal(t, val, "after")

	m := map[string]string{
		"b4": "b4",
	}

	err = set.Field(context.Background(), &m, "dive,keys,default,endkeys,default")
	Equal(t, err, nil)

	val = m["after"]
	Equal(t, val, "after")

	err = set.Field(context.Background(), &m, "keys,endkeys,default")
	Equal(t, err, ErrInvalidKeysTag)

	err = set.Field(context.Background(), &m, "dive,endkeys,default")
	Equal(t, err, ErrUndefinedKeysTag)

	err = set.Field(context.Background(), &m, "dive,keys,undefinedtag")
	Equal(t, err, ErrUndefinedTag{tag: "undefinedtag"})

	err = set.Field(context.Background(), &m, "dive,keys,err,endkeys")
	NotEqual(t, err, nil)

	m = map[string]string{
		"b4": "b4",
	}

	err = set.Field(context.Background(), &m, "dive,keys,default,endkeys,err")
	NotEqual(t, err, nil)
}

func TestStructArray(t *testing.T) {
	type InnerStruct struct {
		String string `s:"defaultStr"`
	}

	type Test struct {
		Inner    InnerStruct
		Arr      []InnerStruct `s:"defaultArr"`
		ArrDive  []InnerStruct `s:"defaultArr,dive"`
		ArrNoTag []InnerStruct
	}

	set := New()
	set.SetTagName("s")
	set.Register("defaultArr", func(ctx context.Context, fl FieldLevel) error {
		if HasValue(fl.Field()) {
			return nil
		}
		fl.Field().Set(reflect.MakeSlice(fl.Field().Type(), 2, 2))
		return nil
	})
	set.Register("defaultStr", func(ctx context.Context, fl FieldLevel) error {
		if fl.Field().String() == "ok" {
			return errors.New("ALREADY OK")
		}
		fl.Field().SetString("default")
		return nil
	})

	var tt Test
	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, len(tt.Arr), 2)
	Equal(t, len(tt.ArrDive), 2)
	Equal(t, tt.Arr[0].String, "")
	Equal(t, tt.Arr[1].String, "")
	Equal(t, tt.ArrDive[0].String, "default")
	Equal(t, tt.ArrDive[1].String, "default")
	Equal(t, tt.Inner.String, "default")

	tt2 := Test{
		Arr: make([]InnerStruct, 1),
	}

	err = set.Struct(context.Background(), &tt2)
	Equal(t, err, nil)
	Equal(t, len(tt2.Arr), 1)
	Equal(t, tt2.Arr[0].String, "")

	tt3 := Test{
		Arr: []InnerStruct{{"ok"}},
	}

	err = set.Struct(context.Background(), &tt3)
	Equal(t, err, nil)
	Equal(t, len(tt3.Arr), 1)
	Equal(t, tt3.Arr[0].String, "ok")

	tt4 := Test{
		ArrDive: []InnerStruct{{"ok"}},
	}

	err = set.Struct(context.Background(), &tt4)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "ALREADY OK")

	tt5 := Test{
		ArrNoTag: make([]InnerStruct, 1),
	}

	err = set.Struct(context.Background(), &tt5)
	Equal(t, err, nil)
	Equal(t, len(tt5.ArrNoTag), 1)
	Equal(t, tt5.ArrNoTag[0].String, "")
}

func TestStructLevel(t *testing.T) {
	type Test struct {
		String string
	}

	set := New()
	set.RegisterStructLevel(func(ctx context.Context, sl StructLevel) error {
		s := sl.Struct().Interface().(Test)
		if s.String == "error" {
			return errors.New("BAD VALUE")
		}
		s.String = "test"
		sl.Struct().Set(reflect.ValueOf(s))
		return nil
	}, Test{})

	var tt Test
	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, tt.String, "test")

	tt.String = "error"
	err = set.Struct(context.Background(), &tt)
	NotEqual(t, err, nil)
}

func TestAlias(t *testing.T) {
	type Test struct {
		String string `r:"repl,repl2"`
	}

	var tt Test
	set := New()
	set.SetTagName("r")
	set.Register("repl", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().SetString("test")
		return nil
	})
	set.Register("repl2", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().SetString("test2")
		return nil
	})

	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, tt.String, "test2")

	set.RegisterAlias("rep", "repl,repl2")
	set.RegisterAlias("bad", "repl,,repl2")
	type Test2 struct {
		String string `r:"rep"`
	}

	var tt2 Test2
	err = set.Struct(context.Background(), &tt2)
	Equal(t, err, nil)
	Equal(t, tt.String, "test2")

	var s string
	err = set.Field(context.Background(), &s, "bad")
	NotEqual(t, err, nil)

	// var s string
	err = set.Field(context.Background(), &s, "repl,rep,bad")
	NotEqual(t, err, nil)
}

func TestParam(t *testing.T) {
	type Test struct {
		String string `r:"ltrim=#$_"`
	}

	set := New()
	set.SetTagName("r")
	set.Register("ltrim", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().SetString(strings.TrimLeft(fl.Field().String(), fl.Param()))
		return nil
	})

	tt := Test{String: "_test"}
	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, tt.String, "test")
}

func TestInterface(t *testing.T) {
	type Test struct {
		Iface interface{} `s:"default"`
	}

	type Inner struct {
		STR    string
		String string `s:"defaultStr"`
	}

	type Test2 struct {
		Iface interface{} `s:"default2"`
	}

	type Inner2 struct {
		String string `s:"error"`
	}

	set := New()
	set.SetTagName("s")
	set.Register("default", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().Set(reflect.ValueOf(Inner{STR: "test"}))
		return nil
	})
	set.Register("default2", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().Set(reflect.ValueOf(Inner2{}))
		return nil
	})
	set.Register("defaultStr", func(ctx context.Context, fl FieldLevel) error {
		if HasValue(fl.Field()) && fl.Field().String() == "ok" {
			return errors.New("ALREADY OK")
		}
		fl.Field().Set(reflect.ValueOf("default"))
		return nil
	})
	set.Register("error", func(ctx context.Context, fl FieldLevel) error {
		return errors.New("BAD VALUE")
	})

	var tt Test
	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	NotEqual(t, tt.Iface, nil)

	inner, ok := tt.Iface.(Inner)
	Equal(t, ok, true)
	Equal(t, inner.String, "default")
	Equal(t, inner.STR, "test")

	var tt2 Test2
	err = set.Struct(context.Background(), &tt2)
	NotEqual(t, err, nil)

	type Test3 struct {
		Iface interface{} `s:"defaultStr"`
	}

	var tt3 Test3
	tt3.Iface = "String"
	err = set.Struct(context.Background(), &tt3)
	Equal(t, err, nil)
	Equal(t, tt3.Iface.(string), "default")

	type Test4 struct {
		Iface interface{} `s:"defaultStr,defaultStr"`
	}

	var tt4 Test4
	tt4.Iface = nil
	err = set.Struct(context.Background(), &tt4)
	Equal(t, err, nil)
	Equal(t, tt4.Iface.(string), "default")

	type Test5 struct {
		Iface interface{} `s:"defaultStr,error"`
	}

	var tt5 Test5
	tt5.Iface = "String"
	err = set.Struct(context.Background(), &tt5)
	NotEqual(t, err, nil)
}

func TestInterfacePtr(t *testing.T) {
	type Test struct {
		Iface interface{} `s:"default"`
	}

	type Inner struct {
		String string `s:"defaultStr"`
	}

	set := New()
	set.SetTagName("s")
	set.Register("default", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().Set(reflect.ValueOf(new(Inner)))
		return nil
	})
	set.Register("defaultStr", func(ctx context.Context, fl FieldLevel) error {
		if fl.Field().String() == "ok" {
			return errors.New("ALREADY OK")
		}
		fl.Field().SetString("default")
		return nil
	})

	var tt Test
	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	NotEqual(t, tt.Iface, nil)

	inner, ok := tt.Iface.(*Inner)
	Equal(t, ok, true)
	Equal(t, inner.String, "default")

	type Test2 struct {
		Iface interface{}
	}

	var tt2 Test2
	tt2.Iface = Inner{}
	err = set.Struct(context.Background(), &tt2)
	Equal(t, err, nil)
}

func TestArray(t *testing.T) {
	type Test struct {
		Arr []string `s:"defaultArr,dive,defaultStr"`
	}

	set := New()
	set.SetTagName("s")
	set.Register("defaultArr", func(ctx context.Context, fl FieldLevel) error {
		if HasValue(fl.Field()) {
			return nil
		}
		fl.Field().Set(reflect.MakeSlice(fl.Field().Type(), 2, 2))
		return nil
	})
	set.Register("defaultStr", func(ctx context.Context, fl FieldLevel) error {
		if fl.Field().String() == "ok" {
			return errors.New("ALREADY OK")
		}
		fl.Field().SetString("default")
		return nil
	})

	var tt Test
	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, len(tt.Arr), 2)
	Equal(t, tt.Arr[0], "default")
	Equal(t, tt.Arr[1], "default")

	tt2 := Test{
		Arr: make([]string, 1),
	}

	err = set.Struct(context.Background(), &tt2)
	Equal(t, err, nil)
	Equal(t, len(tt2.Arr), 1)
	Equal(t, tt2.Arr[0], "default")

	tt3 := Test{
		Arr: []string{"ok"},
	}

	err = set.Struct(context.Background(), &tt3)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "ALREADY OK")
}

func TestMap(t *testing.T) {
	type Test struct {
		Map map[string]string `s:"defaultMap,dive,defaultStr"`
	}

	set := New()
	set.SetTagName("s")
	set.Register("defaultMap", func(ctx context.Context, fl FieldLevel) error {
		if HasValue(fl.Field()) {
			return nil
		}
		fl.Field().Set(reflect.MakeMap(fl.Field().Type()))
		return nil
	})
	set.Register("defaultStr", func(ctx context.Context, fl FieldLevel) error {
		if fl.Field().String() == "ok" {
			return errors.New("ALREADY OK")
		}
		fl.Field().SetString("default")
		return nil
	})

	var tt Test
	err := set.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, len(tt.Map), 0)

	tt2 := Test{
		Map: map[string]string{"key": ""},
	}

	err = set.Struct(context.Background(), &tt2)
	Equal(t, err, nil)
	Equal(t, len(tt2.Map), 1)
	Equal(t, tt2.Map["key"], "default")

	tt3 := Test{
		Map: map[string]string{"key": "ok"},
	}

	err = set.Struct(context.Background(), &tt3)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "ALREADY OK")
}

func TestTimeType(t *testing.T) {
	var tt time.Time
	set := New()
	set.Register("default", func(ctx context.Context, fl FieldLevel) error {
		fl.Field().Set(reflect.ValueOf(time.Now()))
		return nil
	})

	err := set.Field(context.Background(), &tt, "default")
	Equal(t, err, nil)

	err = set.Field(context.Background(), &tt, "default,dive")
	NotEqual(t, err, nil)
	Equal(t, errors.Is(err, ErrInvalidDive), true)
}
