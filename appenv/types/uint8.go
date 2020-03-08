package types

import (
	"strconv"
)

type Uint8PropertyBase struct {
	value uint8
}

func (o *Uint8PropertyBase) Value() interface{} {
	return o.value
}

func (o *Uint8PropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatUint(uint64(o.value), 10)), nil
}

func (o *Uint8PropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseUint(string(text), 10, 8)
	if err != nil {
		return err
	}
	o.value = uint8(v)
	return nil
}

func (o *Uint8PropertyBase) Default() interface{} {
	return uint8(0)
}

var _ Value = (*Uint8PropertyBase)(nil)
