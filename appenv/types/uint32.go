package types

import (
	"strconv"
)

type Uint32PropertyBase struct {
	value uint32
}

func (o *Uint32PropertyBase) Value() interface{} {
	return o.value
}

func (o *Uint32PropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatUint(uint64(o.value), 10)), nil
}

func (o *Uint32PropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseUint(string(text), 10, 32)
	if err != nil {
		return err
	}
	o.value = uint32(v)
	return nil
}

func (o *Uint32PropertyBase) Default() interface{} {
	return uint32(0)
}

var _ Value = (*Uint32PropertyBase)(nil)
