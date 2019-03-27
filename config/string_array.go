package config

import "strings"

type StringArrayConfig []string

// Decode implements the interface `envdecode.Decoder`
func (c *StringArrayConfig) Decode(repl string) error {
	*c = strings.Split(repl, ":")
	return nil
}

// MarshalYAML implements the interface `yaml.Marshaler`
func (c *StringArrayConfig) MarshalYAML() (interface{}, error) {
	return []string(*c), nil
}

// UnmarshalYAML implements the interface `yaml.Unmarshaler`
func (c *StringArrayConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed []string
	if err := unmarshal(&parsed); err != nil {
		return err
	}
	*c = parsed
	return nil
}
