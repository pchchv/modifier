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

## Scrubbers
These functions obfuscate the specified types within the data for pii purposes.

| Name   | Description                                                       |
|--------|-------------------------------------------------------------------|
| emails | Scrubs multiple emails from data.                                 |
| email  | Scrubs the data from and specifies the sha name of the same name. |
| text   | Scrubs the data from and specifies the sha name of the same name. |
| name   | Scrubs the data from and specifies the sha name of the same name. |
| fname  | Scrubs the data from and specifies the sha name of the same name. |
| lname  | Scrubs the data from and specifies the sha name of the same name. |

## Special Notes

`default` and `set` modifiers are special in that they can be used to set the value of a field or underlying type information or attributes and both use the same underlying function to set the data.

Setting a Param will have the following special effects on data types where it's not just the value being set:
- Chan - param used to set the buffer size, default = 0.
- Slice - param used to set the capacity, default = 0.
- Map - param used to set the size, default = 0.
- time.Time - param used to set the time format OR value, default = time.Now(), `utc` = time.Now().UTC(), other tries to parse using RFC3339Nano and set a time value.

To use a comma(,) within your params replace use it's hex representation instead '0x2C' which will be replaced while caching.