# modifier [![GoDoc](https://godoc.org/github.com/pchchv/modifier?status.svg)](https://pkg.go.dev/github.com/pchchv/modifier)

Package modifier is a general library that helps modify or set data in data structures and other objects.

## Installation

#### Use go get

```sh
	go get github.com/pchchv/modifier
```
And

```go
	import "github.com/pchchv/modifier"
```

## Modifiers

These functions modify the data in-place.

| Name                | Description                                                                               |
|---------------------|-------------------------------------------------------------------------------------------|
| camel               | Camel Cases the data.                                                                     |
| default             | Sets the provided default value only if the data is equal to it's default datatype value. |
| empty               | Sets the field equal to the datatype default value. e.g. 0 for int.                       |
| lcase               | lowercases the data.                                                                      |
| ltrim               | Trims spaces from the left of the data provided in the params.                            |
| rtrim               | Trims spaces from the right of the data provided in the params.                           |
| set                 | Set the provided value.                                                                   |
| slug                | Converts the field to a [slug](https://github.com/gosimple/slug)                          |
| snake               | Snake Cases the data.                                                                     |
| strip_alpha         | Strips all ascii characters from the data.                                                |
| strip_alpha_unicode | Strips all unicode characters from the data.                                              |
| strip_num           | Strips all ascii numeric characters from the data.                                        |
| strip_num_unicode   | Strips all unicode numeric characters from the data.                                      |
| strip_punctuation   | Strips all ascii punctuation from the data.                                               |
| title               | Title Cases the data.                                                                     |
| tprefix             | Trims a prefix from the value using the provided param value.                             |
| trim                | Trims space from the data.                                                                |
| tsuffix             | Trims a suffix from the value using the provided param value.                             |
| ucase               | Uppercases the data.                                                                      |
| ucfirst             | Upper cases the first character of the data.                                              |