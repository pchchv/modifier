package main

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
