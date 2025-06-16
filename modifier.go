package modifier

import (
	"context"
	"reflect"
)

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
