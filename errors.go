package modifier

// ErrInvalidTag defines a bad value for a tag being used.
type ErrInvalidTag struct {
	tag   string
	field string
}
