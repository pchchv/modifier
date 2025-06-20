package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pchchv/modifier"
)

var tform *modifier.Transformer

func main() {
	type Test struct {
		String string `modifier:"set"`
	}

	var tt Test
	tform = modifier.New()
	tform.Register("set", transformData)
	if err := tform.Struct(context.Background(), &tt); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", tt)

	var str string
	if err := tform.Field(context.Background(), &str, "set"); err != nil {
		log.Fatal(err)
	}

	fmt.Println(str)
}

func transformData(ctx context.Context, fl modifier.FieldLevel) error {
	fl.Field().SetString("test")
	return nil
}
