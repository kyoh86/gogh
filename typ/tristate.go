package typ

import "fmt"

// Tristate represents a filter state for boolean repository attributes
type Tristate int

const (
	// TristateZero indicates no filtering should be applied
	TristateZero Tristate = iota

	// TristateTrue filters for repositories where the attribute is true
	TristateTrue
	// TristateFalse filters for repositories where the attribute is false
	TristateFalse
)

// AsBoolPtr converts the BooleanFilter to a pointer to a boolean value
func (f Tristate) AsBoolPtr() (*bool, error) {
	var r *bool
	if err := Remap(&r, map[Tristate]*bool{
		TristateZero:  nil,
		TristateTrue:  Ptr(true),
		TristateFalse: Ptr(false),
	}, f); err != nil {
		return nil, fmt.Errorf("invalid Tristate: %w", err)
	}
	return r, nil
}
