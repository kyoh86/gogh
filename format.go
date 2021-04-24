package gogh

import (
	"encoding/json"
	"strings"
)

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

func FormatJSON(p Project) (string, error) {
	buf, err := json.Marshal(map[string]interface{}{
		"fullFilePath": p.FullFilePath(),
		"relFilePath":  p.RelFilePath(),
		"url":          p.URL(),
		"relPath":      p.RelPath(),
		"host":         p.Host(),
		"owner":        p.Owner(),
		"name":         p.Name(),
	})
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func FormatFields(s string) Format {
	return func(p Project) (string, error) {
		return strings.Join([]string{
			p.FullFilePath(),
			p.RelFilePath(),
			p.URL(),
			p.RelPath(),
			p.Host(),
			p.Owner(),
			p.Name(),
		}, s), nil
	}
}
