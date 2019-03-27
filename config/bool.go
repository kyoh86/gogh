package config

import "strconv"

type BoolConfig struct {
	filled bool
	value  bool
}

var (
	TrueConfig = BoolConfig{
		filled: true,
		value:  true,
	}
	FalseConfig = BoolConfig{
		filled: true,
		value:  false,
	}
)

func (c BoolConfig) Bool() bool {
	return c.filled && c.value
}

// Decode implements the interface `envdecode.Decoder`
func (c *BoolConfig) Decode(repl string) error {
	parsed, err := strconv.ParseBool(repl)
	if err != nil {
		return err
	}
	c.filled = true
	c.value = parsed
	return nil
}

// MarshalYAML implements the interface `yaml.Marshaler`
func (c *BoolConfig) MarshalYAML() (interface{}, error) {
	if !c.filled {
		return nil, nil
	}
	return c.value, nil
}

// UnmarshalYAML implements the interface `yaml.Unmarshaler`
func (c *BoolConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed bool
	if err := unmarshal(&parsed); err != nil {
		return err
	}
	c.filled = true
	c.value = parsed
	return nil
}
