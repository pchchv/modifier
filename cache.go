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

type structCache struct {
	lock sync.Mutex
	m    atomic.Value // map[reflect.Type]*cStruct
}
