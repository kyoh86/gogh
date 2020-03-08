package types

import (
	"strconv"
)

type Uint64PropertyBase struct {
	value uint64
}

func (o *Uint64PropertyBase) Value() interface{} {
	return o.value
}

func (o *Uint64PropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatUint(o.value, 10)), nil
}

func (o *Uint64PropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseUint(string(text), 10, 64)
	if err != nil {
		return err
	}
	o.value = uint64(v)
	return nil
}

func (o *Uint64PropertyBase) Default() interface{} {
	return uint64(0)
}

var _ Value = (*Uint64PropertyBase)(nil)
