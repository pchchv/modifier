package modifier

import (
	"sync"
	"sync/atomic"
)

const (
	typeDefault tagType = iota
	typeDive
	typeKeys
	typeEndKeys
)

type tagType uint8

type cTag struct {
	tag            string
	param          string
	aliasTag       string
	actualAliasTag string
	hasAlias       bool
	hasTag         bool
	fn             Func
	keys           *cTag
	next           *cTag
	typeof         tagType
}

type structCache struct {
	lock sync.Mutex
	m    atomic.Value // map[reflect.Type]*cStruct
}
