package modifier

import "reflect"

// FieldLevel represents the interface for field level modifier function.
type FieldLevel interface {
	// Transformer represents a subset of the current *Transformer that is executing the current transformation.
	Transformer() Transform
	// Parent returns the top level parent of the current value return by Field()
	// This is used primarily for having the ability to nil out pointer type values.
	// NOTE: that is there are several layers of abstractions
	//       e. g. interface{} of interface{} of interface{} this
	//       function returns the first interface{}
	Parent() reflect.Value
	// Field returns the current field value being modified.
	Field() reflect.Value
	// Param returns the param associated wth the given function modifier.
	Param() string
}
