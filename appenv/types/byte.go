package types

import (
	"strconv"
)

type BytePropertyBase struct {
	value byte
}

func (o *BytePropertyBase) Value() interface{} {
	return o.value
}

func (o *BytePropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatInt(int64(o.value), 10)), nil
}

func (o *BytePropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseInt(string(text), 10, 8)
	if err != nil {
		return err
	}
	o.value = byte(v)
	return nil
}

func (o *BytePropertyBase) Default() interface{} {
	return byte(0)
}

var _ Value = (*BytePropertyBase)(nil)
