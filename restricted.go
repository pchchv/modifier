package modifier

const (
	diveTag            = "dive"
	keysTag            = "keys"
	ignoreTag          = "-"
	endKeysTag         = "endkeys"
	utf8HexComma       = "0x2C"
	tagSeparator       = ","
	tagKeySeparator    = "="
	restrictedTagChars = ".[],|=+()`~!@#$%^&*\\\"/?<>{}"
)

var (
	restrictedTags = map[string]struct{}{
		diveTag:   {},
		ignoreTag: {},
	}
)
