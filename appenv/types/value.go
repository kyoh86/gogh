package types

import (
	"encoding"
)

type Value interface {
	Value() interface{}
	Default() interface{}
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type Mask interface {
	// Mask secure value.
	Mask(value string) string
}
