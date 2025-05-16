package typ

import "errors"

func Remap[S comparable, V any](v *V, m map[S]V, s S) error {
	var es S
	if s == es {
		return nil
	}
	x, exists := m[s]
	if !exists {
		return errors.New("invalid value")
	}
	*v = x
	return nil
}
