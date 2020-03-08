package types

import (
	"strconv"
)

type Int32PropertyBase struct {
	value int32
}

func (o *Int32PropertyBase) Value() interface{} {
	return o.value
}

func (o *Int32PropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatInt(int64(o.value), 10)), nil
}

func (o *Int32PropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseInt(string(text), 10, 32)
	if err != nil {
		return err
	}
	o.value = int32(v)
	return nil
}

func (o *Int32PropertyBase) Default() interface{} {
	return int32(0)
}

var _ Value = (*Int32PropertyBase)(nil)
