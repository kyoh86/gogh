package types

import (
	"strconv"
)

type Float64PropertyBase struct {
	value float64
}

func (o *Float64PropertyBase) Value() interface{} {
	return o.value
}

func (o *Float64PropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatFloat(float64(o.value), 'f', -1, 64)), nil
}

func (o *Float64PropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}
	o.value = float64(v)
	return nil
}

func (o *Float64PropertyBase) Default() interface{} {
	return float64(0)
}

var _ Value = (*Float64PropertyBase)(nil)
