package prop

import (
	"fmt"
	"reflect"

	"github.com/kyoh86/gogh/appenv/types"
	"github.com/stoewer/go-strcase"
)

type Accessor interface {
	Get() (string, error)
	Set(value string) error
	Unset()
}

type Property struct {
	Type      reflect.Type
	Name      string
	CamelName string
	SnakeName string
	KebabName string

	StoreFile    bool
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

func StoreFile() Store    { return func(d *Property) { d.StoreFile = true } }
func StoreEnvar() Store   { return func(d *Property) { d.StoreEnvar = true } }
func StoreKeyring() Store { return func(d *Property) { d.StoreKeyring = true } }

func Prop(value types.Value, s Store, stores ...Store) (d *Property) {
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
