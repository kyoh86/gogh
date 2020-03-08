package types

type StringPropertyBase struct {
	value string
}

func (o *StringPropertyBase) Value() interface{} {
	return o.value
}

func (o *StringPropertyBase) MarshalText() (text []byte, err error) {
	return []byte(o.value), nil
}

func (o *StringPropertyBase) UnmarshalText(text []byte) error {
	o.value = string(text)
	return nil
}

func (o *StringPropertyBase) Default() interface{} {
	return ""
}

var _ Value = (*StringPropertyBase)(nil)
