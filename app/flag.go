package app

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/pflag"
)

type ProjectFormat string

var _ pflag.Value = (*ProjectFormat)(nil)

func (f ProjectFormat) String() string {
	return string(f)
}

func (f *ProjectFormat) Set(v string) error {
	_, err := formatter(v)
	if err != nil {
		return fmt.Errorf("parse project format: %w", err)
	}
	*f = ProjectFormat(v)
	return nil
}

func (f ProjectFormat) Type() string {
	return "string"
}

func (f ProjectFormat) Formatter() (gogh.Format, error) {
	return formatter(string(f))
}

func formatter(v string) (gogh.Format, error) {
	switch v {
	case "", "rel-path":
		return gogh.FormatRelPath, nil
	case "rel-file-path":
		return gogh.FormatRelFilePath, nil
	case "full-file-path":
		return gogh.FormatFullFilePath, nil
	case "url":
		return gogh.FormatURL, nil
	case "fields":
		return gogh.FormatFields("\t"), nil
	}
	if strings.HasPrefix(v, "fields:") {
		return gogh.FormatFields(v[len("fields:"):]), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}
