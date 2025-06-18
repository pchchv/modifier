package modifier

import (
	"context"
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
