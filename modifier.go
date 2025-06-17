package modifier

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	timeType           = reflect.TypeOf(time.Time{})
	restrictedTagErr   = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
	restrictedAliasErr = "Alias '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
)

// InterceptorFunc is a way to intercept custom types to redirect the functions to be applied to an inner typ/value.
// E. g. sql.NullString, the manipulation should be done on the inner string.
type InterceptorFunc func(current reflect.Value) (inner reflect.Value)

// Func defines a transform function for use.
type Func func(ctx context.Context, fl FieldLevel) error

// StructLevelFunc accepts all values needed for struct level manipulation.
// This is needed for structs that may not be accessed or allowed to add tags from other packages in use.
type StructLevelFunc func(ctx context.Context, sl StructLevel) error

// Transform represents a subset of the
// current *Transformer that is executing the
// current transformation.
type Transform interface {
	Struct(ctx context.Context, v interface{}) error
	Field(ctx context.Context, v interface{}, tags string) error
}

// Transformer is the base controlling object which contains all necessary information
type Transformer struct {
	tagName          string
	aliases          map[string]string
	transformations  map[string]Func
	structLevelFuncs map[reflect.Type]StructLevelFunc
	interceptors     map[reflect.Type]InterceptorFunc
	cCache           *structCache
	tCache           *tagCache
}

// New creates a new Transform object with default tag name of 'mold'.
func New() *Transformer {
	tc := new(tagCache)
	tc.m.Store(make(map[string]*cTag))
	sc := new(structCache)
	sc.m.Store(make(map[reflect.Type]*cStruct))

	return &Transformer{
		tagName:         "mold",
		aliases:         make(map[string]string),
		transformations: make(map[string]Func),
		interceptors:    make(map[reflect.Type]InterceptorFunc),
		cCache:          sc,
		tCache:          tc,
	}
}

// // SetTagName sets the given tag name to be used.
// // Default is "trans"
// func (t *Transformer) SetTagName(tagName string) {
// 	t.tagName = tagName
// }

// Register adds a transformation with the given tag.
//
// NOTES:
// - if the key already exists, the previous transformation function will be replaced.
// - this method is not thread-safe it is intended that these all be registered before hand.
func (t *Transformer) Register(tag string, fn Func) {
	if len(tag) == 0 {
		panic("Function Key cannot be empty")
	}

	if fn == nil {
		panic("Function cannot be empty")
	}

	if _, ok := restrictedTags[tag]; ok || strings.ContainsAny(tag, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedTagErr, tag))
	}

	t.transformations[tag] = fn
}

// RegisterStructLevel registers StructLevelFunc against a number of types.
// This is needed for structs that may not be access or rights to add tags from other packages in use.
//
// NOTE: this method is not thread-safe. It is intended that all of them must be registered prior to any validation.
func (t *Transformer) RegisterStructLevel(fn StructLevelFunc, types ...interface{}) {
	if t.structLevelFuncs == nil {
		t.structLevelFuncs = make(map[reflect.Type]StructLevelFunc)
	}

	for _, typ := range types {
		t.structLevelFuncs[reflect.TypeOf(typ)] = fn
	}
}

// RegisterAlias registers a mapping of a single transform tag that defines a common or
// complex set of transformations to simplify adding transforms to structs.
//
// NOTE: this method is not thread-safe. It is intended that these all be registered before hand.
func (t *Transformer) RegisterAlias(alias, tags string) {
	if len(alias) == 0 {
		panic("Alias cannot be empty")
	}

	if len(tags) == 0 {
		panic("Aliased tags cannot be empty")
	}

	if _, ok := restrictedTags[alias]; ok || strings.ContainsAny(alias, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedAliasErr, alias))
	}

	t.aliases[alias] = tags
}

// RegisterInterceptor registers a new interceptor functions agains one or more types.
// This InterceptorFunc allows one to intercept the incoming to redirect the
// application of modifications to an inner type/value.
// E. g. sql.NullString
func (t *Transformer) RegisterInterceptor(fn InterceptorFunc, types ...interface{}) {
	for _, typ := range types {
		t.interceptors[reflect.TypeOf(typ)] = fn
	}
}

func (t *Transformer) setByField(ctx context.Context, orig reflect.Value, ct *cTag) (err error) {
	current, kind := t.extractType(orig)
	if ct != nil && ct.hasTag {
		for ct != nil {
			switch ct.typeof {
			case typeEndKeys:
				return
			case typeDive:
				ct = ct.next
				switch kind {
				case reflect.Slice, reflect.Array:
					err = t.setByIterable(ctx, current, ct)
				case reflect.Map:
					err = t.setByMap(ctx, current, ct)
				case reflect.Ptr:
					innerKind := current.Type().Elem().Kind()
					if innerKind == reflect.Slice || innerKind == reflect.Map {
						// is a nil pointer to a slice or map, nothing to do.
						return nil
					}
					// not a valid use of the dive tag
					fallthrough
				default:
					err = ErrInvalidDive
				}
				return
			default:
				if !current.CanAddr() {
					newVal := reflect.New(current.Type()).Elem()
					newVal.Set(current)
					if err = ct.fn(ctx, fieldLevel{
						transformer: t,
						parent:      orig,
						current:     newVal,
						param:       ct.param,
					}); err != nil {
						return
					}
					orig.Set(reflect.Indirect(newVal))
					current, kind = t.extractType(orig)
				} else {
					if err = ct.fn(ctx, fieldLevel{
						transformer: t,
						parent:      orig,
						current:     current,
						param:       ct.param,
					}); err != nil {
						return
					}
					// value could have been changed or reassigned
					current, kind = t.extractType(current)
				}
				ct = ct.next
			}
		}
	}

	// need to do this again because one of the
	// previous sets could have set a struct value,
	// where it was a nil pointer before
	orig2 := current
	current, kind = t.extractType(current)
	if kind == reflect.Struct {
		typ := current.Type()
		if typ == timeType {
			return
		}

		if !current.CanAddr() {
			newVal := reflect.New(typ).Elem()
			newVal.Set(current)

			if err = t.setByStruct(ctx, orig, newVal, typ); err != nil {
				return
			}
			orig.Set(reflect.Indirect(newVal))
			return
		}
		err = t.setByStruct(ctx, orig2, current, typ)
	}
	return
}

func (t *Transformer) setByMap(ctx context.Context, current reflect.Value, ct *cTag) error {
	for _, key := range current.MapKeys() {
		newVal := reflect.New(current.Type().Elem()).Elem()
		newVal.Set(current.MapIndex(key))
		if ct != nil && ct.typeof == typeKeys && ct.keys != nil {
			// remove current map key as we may be changing it
			// and re-add to the map afterwards
			current.SetMapIndex(key, reflect.Value{})
			newKey := reflect.New(current.Type().Key()).Elem()
			newKey.Set(key)
			key = newKey
			// handle map key
			if err := t.setByField(ctx, key, ct.keys); err != nil {
				return err
			}

			// can be nil when just keys being validated
			if ct.next != nil {
				if err := t.setByField(ctx, newVal, ct.next); err != nil {
					return err
				}
			}
		} else {
			if err := t.setByField(ctx, newVal, ct); err != nil {
				return err
			}
		}
		current.SetMapIndex(key, newVal)
	}

	return nil
}
