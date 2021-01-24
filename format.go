package gogh

import "strings"

type Format func(Project) (string, error)

func (f Format) Format(p Project) (string, error) {
	return f(p)
}

func FormatFullFilePath(p Project) (string, error) {
	return p.FullFilePath(), nil
}

func FormatRelPath(p Project) (string, error) {
	return p.RelPath(), nil
}

func FormatRelFilePath(p Project) (string, error) {
	return p.RelFilePath(), nil
}

func FormatURL(p Project) (string, error) {
	return p.URL(), nil
}

func FormatFields(s string) Format {
	return func(p Project) (string, error) {
		return strings.Join([]string{
			p.FullFilePath(),
			p.RelPath(),
			p.Host(),
			p.User(),
			p.Name(),
		}, s), nil
	}
}
