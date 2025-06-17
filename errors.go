package modifier

import (
	"fmt"
	"strings"
)

// ErrInvalidTag defines a bad value for a tag being used.
type ErrInvalidTag struct {
	tag   string
	field string
}

// Error returns the InvalidTag error text.
func (e *ErrInvalidTag) Error() string {
	return fmt.Sprintf("invalid tag '%s' found on field %s", e.tag, e.field)
}

// ErrUndefinedTag defines a tag that does not exist.
type ErrUndefinedTag struct {
	tag   string
	field string
}

// Error returns the UndefinedTag error text.
func (e *ErrUndefinedTag) Error() string {
	return strings.TrimSpace(fmt.Sprintf("unregistered/undefined transformation '%s' found on field %s", e.tag, e.field))
}
