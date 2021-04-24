package view

import (
	"encoding/json"
	"strings"

	"github.com/kyoh86/gogh/v2"
)

type ProjectFormat func(gogh.Project) (string, error)

func (f ProjectFormat) Format(p gogh.Project) (string, error) {
	return f(p)
}

func ProjectFormatFullFilePath(p gogh.Project) (string, error) {
	return p.FullFilePath(), nil
}

func ProjectFormatRelPath(p gogh.Project) (string, error) {
	return p.RelPath(), nil
}

func ProjectFormatRelFilePath(p gogh.Project) (string, error) {
	return p.RelFilePath(), nil
}

func ProjectFormatURL(p gogh.Project) (string, error) {
	return p.URL(), nil
}

func ProjectFormatJSON(p gogh.Project) (string, error) {
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

func ProjectFormatFields(s string) ProjectFormat {
	return func(p gogh.Project) (string, error) {
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
