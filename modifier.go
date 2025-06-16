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

// Transform represents a subset of the
// current *Transformer that is executing the
// current transformation.
type Transform interface {
	Struct(ctx context.Context, v interface{}) error
	Field(ctx context.Context, v interface{}, tags string) error
}
