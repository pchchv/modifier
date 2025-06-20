package main

import (
	"context"

	"github.com/pchchv/modifier"
)

func main() {}

func transformData(ctx context.Context, fl modifier.FieldLevel) error {
	fl.Field().SetString("test")
	return nil
}
