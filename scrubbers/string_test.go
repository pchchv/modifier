package scrubbers

import (
	"context"
	"testing"

	. "github.com/pchchv/go-assert"
)

func TestEmails(t *testing.T) {
	type Test struct {
		Email string `scrub:"emails"`
	}

	scrub := New()
	email := "Jack.Pochechuev@gmail.com"
	tt := Test{Email: email}
	err := scrub.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, tt.Email, "<<scrubbed::email::sha1::5131512f2d165ca283b055bc6f32bc01dd23121e>>@gmail.com")

	err = scrub.Field(context.Background(), &email, "emails")
	Equal(t, err, nil)
	Equal(t, email, "<<scrubbed::email::sha1::5131512f2d165ca283b055bc6f32bc01dd23121e>>@gmail.com")

	var iface interface{}
	err = scrub.Field(context.Background(), &iface, "emails")
	Equal(t, err, nil)
	Equal(t, iface, nil)

	iface = "Jack.Pochechuev@gmail.com"
	err = scrub.Field(context.Background(), &iface, "emails")
	Equal(t, err, nil)
	Equal(t, iface, "<<scrubbed::email::sha1::5131512f2d165ca283b055bc6f32bc01dd23121e>>@gmail.com")
}
