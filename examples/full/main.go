package main

import "net/url"

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

func main() {}

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
