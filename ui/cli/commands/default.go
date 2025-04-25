package commands

func DefaultValue[T comparable](value T, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

func DefaultSlice[T any](value, defaultValue []T) []T {
	if value == nil {
		return defaultValue
	}
	return value
}
