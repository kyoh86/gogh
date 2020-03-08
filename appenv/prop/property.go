package prop

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kyoh86/gogh/appenv/internal"
	"github.com/kyoh86/gogh/appenv/types"
	"github.com/stoewer/go-strcase"
)

type Accessor interface {
	Get() (string, error)
	Set(value string) error
	Unset()
}

type Store func(d *internal.Property)

// NOTE: StoreXXX must not be merged to the Value.
// They can be expand for other storages, but that is NOT
// realistic to implement all of them in all of properties.

func File() Store    { return func(d *internal.Property) { d.StoreFile = true } }
func Envar() Store   { return func(d *internal.Property) { d.StoreEnvar = true } }
func Keyring() Store { return func(d *internal.Property) { d.StoreKeyring = true } }

func Prop(value types.Value, s Store, stores ...Store) (d *internal.Property) {
	d = new(internal.Property)
	typ := reflect.ValueOf(value).Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	d.PkgPath = typ.PkgPath()
	d.Name = typ.Name()

	valueType := reflect.ValueOf(value.Value()).Type()
	d.ValueEmpty = reflect.Zero(valueType).Interface()
	d.ValueType = fmt.Sprintf("%T", d.ValueEmpty)

	s(d)
	for _, s := range stores {
		s(d)
	}

	d.CamelName = strcase.LowerCamelCase(d.Name)
	d.SnakeName = strcase.UpperSnakeCase(d.Name)
	d.KebabName = strcase.KebabCase(d.Name)
	d.DottedName = strings.ReplaceAll(d.KebabName, "-", ".")
	return
}

type Mask interface {
	// Mask secure value.
	Mask(value string) string
}
