package props

import (
	"encoding"
	"fmt"
	"reflect"

	"github.com/stoewer/go-strcase"
)

type Value interface {
	Value() interface{}
	Default() interface{}
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type Property struct {
	Type      reflect.Type
	Name      string
	CamelName string
	SnakeName string
	KebabName string

	StoreConfig  bool
	StoreCache   bool
	StoreEnvar   bool
	StoreKeyring bool

	ValueType   reflect.Type
	ValueEmpty  interface{}
	ValueTypeID string
}

type Store func(d *Property)

// NOTE: StoreXXX must not be merged to the Value.
// They can be expand for other storages, but that is NOT
// realistic to implement all of them in all of properties.

func StoreConfig() Store  { return func(d *Property) { d.StoreConfig = true } }
func StoreCache() Store   { return func(d *Property) { d.StoreCache = true } }
func StoreEnvar() Store   { return func(d *Property) { d.StoreEnvar = true } }
func StoreKeyring() Store { return func(d *Property) { d.StoreKeyring = true } }

func Prop(value Value, s Store, stores ...Store) (d *Property) {
	d = new(Property)
	d.Type = reflect.ValueOf(value).Type()
	for d.Type.Kind() == reflect.Ptr {
		d.Type = d.Type.Elem()
	}

	d.Name = d.Type.Name()

	d.ValueType = reflect.ValueOf(value.Value()).Type()
	d.ValueEmpty = reflect.Zero(d.ValueType).Interface()
	d.ValueTypeID = fmt.Sprintf("%T", d.ValueEmpty)

	s(d)
	for _, s := range stores {
		s(d)
	}

	d.CamelName = strcase.LowerCamelCase(d.Name)
	d.SnakeName = strcase.UpperSnakeCase(d.Name)
	d.KebabName = strcase.KebabCase(d.Name)
	return
}

type Mask interface {
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

var _ Value = (*StringPropertyBase)(nil)
var _ encoding.TextMarshaler = (*StringPropertyBase)(nil)
var _ encoding.TextUnmarshaler = (*StringPropertyBase)(nil)
