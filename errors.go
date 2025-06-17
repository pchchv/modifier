package modifier

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidDive describes an invalid dive tag configuration.
	ErrInvalidDive = errors.New("invalid dive tag configuration")
	// ErrInvalidKeysTag describes a misuse of the keys tag.
	ErrInvalidKeysTag = errors.New("'" + keysTag + "' tag must be immediately preceeded by the '" + diveTag + "' tag")
	// ErrUndefinedKeysTag describes an undefined keys tag when and endkeys tag defined.
	ErrUndefinedKeysTag = errors.New("'" + endKeysTag + "' tag encountered without a corresponding '" + keysTag + "' tag")
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
