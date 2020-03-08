package types

import (
	"strconv"
)

type Int64PropertyBase struct {
	value int64
}

func (o *Int64PropertyBase) Value() interface{} {
	return o.value
}

func (o *Int64PropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatInt(int64(o.value), 10)), nil
}

func (o *Int64PropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return err
	}
	o.value = int64(v)
	return nil
}

func (o *Int64PropertyBase) Default() interface{} {
	return int64(0)
}

var _ Value = (*Int64PropertyBase)(nil)
