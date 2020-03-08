package types

import (
	"strconv"
)

type BoolPropertyBase struct {
	value bool
}

func (o *BoolPropertyBase) Value() interface{} {
	return o.value
}

func (o *BoolPropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatBool(o.value)), nil
}

func (o *BoolPropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseBool(string(text))
	if err != nil {
		return err
	}
	o.value = v
	return nil
}

func (o *BoolPropertyBase) Default() interface{} {
	return false
}

var _ Value = (*BoolPropertyBase)(nil)
