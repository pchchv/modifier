package modifier

import "reflect"

// extractType gets the actual underlying type of field value.
func (t *Transformer) extractType(current reflect.Value) (reflect.Value, reflect.Kind) {
	switch current.Kind() {
	case reflect.Ptr:
		if current.IsNil() {
			return current, reflect.Ptr
		} else {
			return t.extractType(current.Elem())
		}
	case reflect.Interface:
		if current.IsNil() {
			return current, reflect.Interface
		} else {
			return t.extractType(current.Elem())
		}
	default:
		if fn := t.interceptors[current.Type()]; fn != nil {
			return t.extractType(fn(current))
		} else {
			return current, current.Kind()
		}
	}
}
