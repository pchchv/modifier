package modifier

import "fmt"

// ErrInvalidTag defines a bad value for a tag being used.
type ErrInvalidTag struct {
	tag   string
	field string
}

// Error returns the InvalidTag error text.
func (e *ErrInvalidTag) Error() string {
	return fmt.Sprintf("invalid tag '%s' found on field %s", e.tag, e.field)
}
