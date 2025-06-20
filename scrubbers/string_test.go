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

func TestText(t *testing.T) {
	type Test struct {
		String string `scrub:"text"`
	}

	scrub := New()
	name := "Joey Bloggs"
	tt := Test{String: name}
	err := scrub.Struct(context.Background(), &tt)
	Equal(t, err, nil)
	Equal(t, tt.String, "<<scrubbed::text::sha1::028f74c1850aa1efb33a2e8746c0f4183e1e8e30>>")

	err = scrub.Field(context.Background(), &name, "text")
	Equal(t, err, nil)
	Equal(t, name, "<<scrubbed::text::sha1::028f74c1850aa1efb33a2e8746c0f4183e1e8e30>>")

	var iface interface{}
	err = scrub.Field(context.Background(), &iface, "text")
	Equal(t, err, nil)
	Equal(t, iface, nil)

	iface = "Joey Bloggs"
	err = scrub.Field(context.Background(), &iface, "text")
	Equal(t, err, nil)
	Equal(t, iface, "<<scrubbed::text::sha1::028f74c1850aa1efb33a2e8746c0f4183e1e8e30>>")

	// testing Text wrapped func
	name = "Joey Bloggs"
	err = scrub.Field(context.Background(), &name, "name")
	Equal(t, err, nil)
	Equal(t, name, "<<scrubbed::name::sha1::028f74c1850aa1efb33a2e8746c0f4183e1e8e30>>")
}
