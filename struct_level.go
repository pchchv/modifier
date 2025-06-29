package modifier

import "reflect"

var _ StructLevel = (*structLevel)(nil)

// StructLevel represents the interface for struct level modifier function
type StructLevel interface {
	// Transformer represents a subset of the current *Transformer that is executing the current transformation.
	Transformer() Transform
	// Parent returns the top level parent of the current value return by Struct().
	// This is used primarily for having the ability to nil out pointer type values.
	// NOTE: that is there are several layers of abstractions
	//       e. g. interface{} of interface{} of interface{} this
	//       function returns the first interface{}.
	Parent() reflect.Value
	// Struct returns the value of the current struct being modified.
	Struct() reflect.Value
}

type structLevel struct {
	transformer *Transformer
	parent      reflect.Value
	current     reflect.Value
}

func (s structLevel) Parent() reflect.Value {
	return s.parent
}

func (s structLevel) Struct() reflect.Value {
	return s.current
}

func (s structLevel) Transformer() Transform {
	return s.transformer
}
