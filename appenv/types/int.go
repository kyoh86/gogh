package types

import (
	"strconv"
)

type IntPropertyBase struct {
	value int
}

func (o *IntPropertyBase) Value() interface{} {
	return o.value
}

func (o *IntPropertyBase) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatInt(int64(o.value), 10)), nil
}

func (o *IntPropertyBase) UnmarshalText(text []byte) error {
	v, err := strconv.ParseInt(string(text), 10, strconv.IntSize)
	if err != nil {
		return err
	}
	o.value = int(v)
	return nil
}

func (o *IntPropertyBase) Default() interface{} {
	return int(0)
}

var _ Value = (*IntPropertyBase)(nil)
