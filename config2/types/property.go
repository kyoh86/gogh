package types

import "encoding"

type Property interface {
	Value() interface{}
	Default() interface{}
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type StoreConfigFile interface{ StoreConfigFile() }
type StoreCacheFile interface{ StoreCacheFile() }
type StoreEnvar interface{ StoreEnvar() }
type StoreKeyring interface{ StoreKeyring() }

type Masking interface {
	// Mask secure value.
	Mask(value string) string
}

type StringPropertyBase struct {
	value string
}

func (o *StringPropertyBase) Value() interface{} {
	return o.value
}

func (o *StringPropertyBase) MarshalText() (text []byte, err error) {
	return []byte(o.value), nil
}

func (o *StringPropertyBase) UnmarshalText(text []byte) error {
	o.value = string(text)
	return nil
}

func (o *StringPropertyBase) Default() interface{} {
	return ""
}

var _ Property = (*StringPropertyBase)(nil)
var _ encoding.TextMarshaler = (*StringPropertyBase)(nil)
var _ encoding.TextUnmarshaler = (*StringPropertyBase)(nil)
