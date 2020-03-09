package gen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kyoh86/gogh/appenv/types"
	"github.com/stoewer/go-strcase"
)

type Store interface {
	mark(d *Property)
}

type storeFunc func(d *Property)

func (f storeFunc) mark(d *Property) {
	f(d)
}

// NOTE: StoreXXX must not be merged to the Value.
// They can be expand for other storages, but that is NOT
// realistic to implement all of them in all of properties.

func YAML() Store    { return storeFunc(func(d *Property) { d.storeYAML = true }) }
func Envar() Store   { return storeFunc(func(d *Property) { d.storeEnvar = true }) }
func Keyring() Store { return storeFunc(func(d *Property) { d.storeKeyring = true }) }

func Prop(value types.Value, s Store, stores ...Store) (d *Property) {
	d = new(Property)
	typ := reflect.ValueOf(value).Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	d.pkgPath = typ.PkgPath()
	d.name = typ.Name()

	valueType := reflect.ValueOf(value.Value()).Type()
	d.valueEmpty = reflect.Zero(valueType).Interface()
	d.valueType = fmt.Sprintf("%T", d.valueEmpty)

	s.mark(d)
	for _, s := range stores {
		s.mark(d)
	}
	_, d.mask = value.(types.Mask)

	d.camelName = strcase.LowerCamelCase(d.name)
	d.snakeName = strcase.UpperSnakeCase(d.name)
	d.kebabName = strcase.KebabCase(d.name)
	d.dottedName = strings.ReplaceAll(d.kebabName, "-", ".")
	return
}

// Property describes environment property.
// It is in internal package, and it can be generated with prop.Prop.
type Property struct {
	pkgPath string
	name    string

	camelName  string
	snakeName  string
	kebabName  string
	dottedName string

	storeYAML    bool
	storeEnvar   bool
	storeKeyring bool

	mask bool

	valueEmpty interface{}
	valueType  string
}
