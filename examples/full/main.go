package main

import (
	"context"
	"log"
	"net/url"

	"github.com/pchchv/form"
	"github.com/pchchv/modifier/modifiers"
	"github.com/pchchv/modifier/scrubbers"
	"github.com/pchchv/validator"
)

var (
	decoder  = form.NewDecoder()
	scrub    = scrubbers.New()
	conform  = modifiers.New()
	validate = validator.New()
)

// Address contains address information.
type Address struct {
	Name  string `mod:"trim" validate:"required"`
	Phone string `mod:"trim" validate:"required"`
}

// User contains user information.
type User struct {
	Name    string            `mod:"trim"      validate:"required"              scrub:"name"`
	Age     uint8             `                validate:"required,gt=0,lt=130"`
	Gender  string            `                validate:"required"`
	Email   string            `mod:"trim"      validate:"required,email"        scrub:"emails"`
	Address []Address         `                validate:"required,dive"`
	Active  bool              `form:"active"`
	Misc    map[string]string `mod:"dive,keys,trim,endkeys,trim"`
}

// parseForm simulates the results of http.Request's ParseForm() function.
func parseForm() url.Values {
	return url.Values{
		"Name":             []string{"  joeybloggs  "},
		"Age":              []string{"3"},
		"Gender":           []string{"Male"},
		"Email":            []string{"Dean.Karn@gmail.com  "},
		"Address[0].Name":  []string{"26 Here Blvd."},
		"Address[0].Phone": []string{"9(999)999-9999"},
		"Address[1].Name":  []string{"26 There Blvd."},
		"Address[1].Phone": []string{"1(111)111-1111"},
		"active":           []string{"true"},
		"Misc[  b4  ]":     []string{"  b4  "},
	}
}

func main() {
	var user User
	values := parseForm()
	// must pass a pointer
	if err := decoder.Decode(&user, values); err != nil {
		log.Panic(err)
	}

	log.Printf("Decoded:%+v\n\n", user)

	// great now lets conform our values,
	// after all a human input the data nobody's perfect
	if err := conform.Struct(context.Background(), &user); err != nil {
		log.Panic(err)
	}

	log.Printf("Conformed:%+v\n\n", user)

	// that's better all those extra spaces are gone let's validate the data
	if err := validate.Struct(user); err != nil {
		log.Panic(err)
	}

	// data's validated, proceeding:
	// save to database
	// process request
	// etc....

	// working with the data is over, let's log it or save it somewhere
	// there is sensitive PII data, it should be made sure that it is de-identified first
	if err := scrub.Struct(context.Background(), &user); err != nil {
		log.Panic(err)
	}

	log.Printf("Scrubbed:%+v\n\n", user)
}
