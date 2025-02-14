package main

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v3/view"
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

func (f ProjectFormat) Formatter() (view.ProjectFormat, error) {
	return formatter(string(f))
}

func formatter(v string) (view.ProjectFormat, error) {
	switch v {
	case "", "rel-path", "rel":
		return view.ProjectFormatRelPath, nil
	case "rel-file-path":
		return view.ProjectFormatRelFilePath, nil
	case "full-file-path", "full":
		return view.ProjectFormatFullFilePath, nil
	case "json":
		return view.ProjectFormatJSON, nil
	case "url":
		return view.ProjectFormatURL, nil
	case "fields":
		return view.ProjectFormatFields("\t"), nil
	}
	if strings.HasPrefix(v, "fields:") {
		return view.ProjectFormatFields(v[len("fields:"):]), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}
