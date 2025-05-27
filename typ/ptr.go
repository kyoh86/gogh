package typ

func NilablePtr[T comparable](v T) *T {
	if v == defaultValue[T]() {
		return nil
	}
	return &v
}

// defaultValue returns the zero value for a given type.
func defaultValue[T any]() T {
	var zero T
	return zero
}

// Ptr converts a value of any type to a pointer to that value.
func Ptr[T any](v T) *T {
	return &v
}
