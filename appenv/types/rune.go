package types

import "errors"

type RunePropertyBase struct {
	value rune
}

func (o *RunePropertyBase) Value() interface{} {
	return o.value
}

func (o *RunePropertyBase) MarshalText() (text []byte, err error) {
	return []byte(string([]rune{o.value})), nil
}

func (o *RunePropertyBase) UnmarshalText(text []byte) error {
	runes := []rune(string(text))
	if len(runes) != 1 {
		return errors.New("invalid rune")
	}
	o.value = runes[0]
	return nil
}

func (o *RunePropertyBase) Default() interface{} {
	return rune(0)
}

var _ Value = (*RunePropertyBase)(nil)
