package modifier

import (
	"reflect"
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

type cField struct {
	idx   int
	cTags *cTag
}

type cStruct struct {
	fields []*cField
	fn     StructLevelFunc
}

type structCache struct {
	lock sync.Mutex
	m    atomic.Value // map[reflect.Type]*cStruct
}

func (sc *structCache) Get(key reflect.Type) (c *cStruct, found bool) {
	c, found = sc.m.Load().(map[reflect.Type]*cStruct)[key]
	return
}

func (sc *structCache) Set(key reflect.Type, value *cStruct) {
	m := sc.m.Load().(map[reflect.Type]*cStruct)
	nm := make(map[reflect.Type]*cStruct, len(m)+1)
	for k, v := range m {
		nm[k] = v
	}

	nm[key] = value
	sc.m.Store(nm)
}

type tagCache struct {
	lock sync.Mutex
	m    atomic.Value // map[string]*cTag
}
