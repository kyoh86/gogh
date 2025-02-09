package gogh

// falsePtr converts bool (default: false) to *bool (default: nil as true)
func falsePtr(b bool) *bool {
	if b {
		f := false
		return &f
	}
	return nil // == &true
}

// nilablePtr converts a value of any type to a pointer, returning nil if the value is the zero value.
func nilablePtr[T comparable](v T) *T {
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

// ptr converts a value of any type to a pointer to that value.
func ptr[T any](v T) *T {
	return &v
}
