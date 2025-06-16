package modifier

import "reflect"

// InterceptorFunc is a way to intercept custom types to redirect the functions to be applied to an inner typ/value.
// E. g. sql.NullString, the manipulation should be done on the inner string.
type InterceptorFunc func(current reflect.Value) (inner reflect.Value)
