package types

import (
	"strconv"
)

type UintPropertyBase struct {
	value uint
}

func (o *UintPropertyBase) Value() interface{} {
	return o.value
}

func (o *UintPropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatUint(uint64(o.value), 10)), nil
}

func (o *UintPropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseUint(string(text), 10, strconv.IntSize)
	if err != nil {
		return err
	}
	o.value = uint(v)
	return nil
}

func (o *UintPropertyBase) Default() interface{} {
	return uint(0)
}

var _ Value = (*UintPropertyBase)(nil)
