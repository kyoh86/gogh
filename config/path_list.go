package config

import (
	"path/filepath"
)

type PathListOption []string

// Decode implements the interface `envdecode.Decoder`
func (c *PathListOption) Decode(repl string) error {
	*c = filepath.SplitList(repl)
	return nil
}

// MarshalYAML implements the interface `yaml.Marshaler`
func (c PathListOption) MarshalYAML() (interface{}, error) {
	return []string(c), nil
}

// UnmarshalYAML implements the interface `yaml.Unmarshaler`
func (c *PathListOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed []string
	if err := unmarshal(&parsed); err != nil {
		return err
	}
	*c = parsed
	return nil
}
