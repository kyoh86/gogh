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

func (f ProjectFormat) Formatter() (gogh.ProjectFormat, error) {
	return formatter(string(f))
}

func formatter(v string) (gogh.ProjectFormat, error) {
	switch v {
	case "", "rel-path":
		return gogh.ProjectFormatRelPath, nil
	case "rel-file-path":
		return gogh.ProjectFormatRelFilePath, nil
	case "full-file-path":
		return gogh.ProjectFormatFullFilePath, nil
	case "json":
		return gogh.ProjectFormatJSON, nil
	case "url":
		return gogh.ProjectFormatURL, nil
	case "fields":
		return gogh.ProjectFormatFields("\t"), nil
	}
	if strings.HasPrefix(v, "fields:") {
		return gogh.ProjectFormatFields(v[len("fields:"):]), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}
