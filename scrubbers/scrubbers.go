package scrubbers

import "github.com/pchchv/modifier"

// New returns a scrubber with defaults registered.
func New() *modifier.Transformer {
	scrub := modifier.New()
	scrub.SetTagName("scrub")
	scrub.Register("emails", emails)
	scrub.Register("text", textFn("text"))
	scrub.Register("email", textFn("email"))
	scrub.Register("name", textFn("name"))
	scrub.Register("fname", textFn("fname"))
	scrub.Register("lname", textFn("lname"))
	return scrub
}
