package internal

// Property describes environment property.
// It is in internal package, and it can be generated with prop.Prop.
type Property struct {
	PkgPath    string
	Name       string
	CamelName  string
	SnakeName  string
	KebabName  string
	DottedName string

	StoreFile    bool
	StoreEnvar   bool
	StoreKeyring bool

	ValueEmpty interface{}
	ValueType  string
}
