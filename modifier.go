package modifier

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

var restrictedTagErr = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"

// InterceptorFunc is a way to intercept custom types to redirect the functions to be applied to an inner typ/value.
// E. g. sql.NullString, the manipulation should be done on the inner string.
type InterceptorFunc func(current reflect.Value) (inner reflect.Value)

// Func defines a transform function for use.
type Func func(ctx context.Context, fl FieldLevel) error

// StructLevelFunc accepts all values needed for struct level manipulation.
// This is needed for structs that may not be accessed or allowed to add tags from other packages in use.
type StructLevelFunc func(ctx context.Context, sl StructLevel) error

// Transform represents a subset of the
// current *Transformer that is executing the
// current transformation.
type Transform interface {
	Struct(ctx context.Context, v interface{}) error
	Field(ctx context.Context, v interface{}, tags string) error
}

// Transformer is the base controlling object which contains all necessary information
type Transformer struct {
	tagName          string
	aliases          map[string]string
	transformations  map[string]Func
	structLevelFuncs map[reflect.Type]StructLevelFunc
	interceptors     map[reflect.Type]InterceptorFunc
	cCache           *structCache
	tCache           *tagCache
}

// New creates a new Transform object with default tag name of 'mold'.
func New() *Transformer {
	tc := new(tagCache)
	tc.m.Store(make(map[string]*cTag))
	sc := new(structCache)
	sc.m.Store(make(map[reflect.Type]*cStruct))

	return &Transformer{
		tagName:         "mold",
		aliases:         make(map[string]string),
		transformations: make(map[string]Func),
		interceptors:    make(map[reflect.Type]InterceptorFunc),
		cCache:          sc,
		tCache:          tc,
	}
}

// // SetTagName sets the given tag name to be used.
// // Default is "trans"
// func (t *Transformer) SetTagName(tagName string) {
// 	t.tagName = tagName
// }

// Register adds a transformation with the given tag.
//
// NOTES:
// - if the key already exists, the previous transformation function will be replaced.
// - this method is not thread-safe it is intended that these all be registered before hand.
func (t *Transformer) Register(tag string, fn Func) {
	if len(tag) == 0 {
		panic("Function Key cannot be empty")
	}

	if fn == nil {
		panic("Function cannot be empty")
	}

	if _, ok := restrictedTags[tag]; ok || strings.ContainsAny(tag, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedTagErr, tag))
	}

	t.transformations[tag] = fn
}

// RegisterStructLevel registers StructLevelFunc against a number of types.
// This is needed for structs that may not be access or rights to add tags from other packages in use.
//
// NOTE: this method is not thread-safe. It is intended that all of them must be registered prior to any validation.
func (t *Transformer) RegisterStructLevel(fn StructLevelFunc, types ...interface{}) {
	if t.structLevelFuncs == nil {
		t.structLevelFuncs = make(map[reflect.Type]StructLevelFunc)
	}

	for _, typ := range types {
		t.structLevelFuncs[reflect.TypeOf(typ)] = fn
	}
}
