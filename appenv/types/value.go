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
