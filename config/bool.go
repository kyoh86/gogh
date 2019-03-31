package config

import (
	"errors"
	"strings"
)

type BoolOption string

var (
	TrueOption      = BoolOption("yes")
	FalseOption     = BoolOption("no")
	EmptyBoolOption = BoolOption("")
)

func (c BoolOption) String() string {
	return string(c)
}

func (c BoolOption) Bool() bool {
	return c == TrueOption
}

// Decode implements the interface `envdecode.Decoder`
func (c *BoolOption) Decode(repl string) error {
	switch strings.ToLower(repl) {
	case "yes", "no", "":
		*c = BoolOption(repl)
		return nil
	}
	return errors.New("invalid type")
}

// MarshalYAML implements the interface `yaml.Marshaler`
func (c BoolOption) MarshalYAML() (interface{}, error) {
	return string(c), nil
}

// UnmarshalYAML implements the interface `yaml.Unmarshaler`
func (c *BoolOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed string
	if err := unmarshal(&parsed); err != nil {
		return err
	}
	return c.Decode(parsed)
}
