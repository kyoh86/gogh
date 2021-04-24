package view

import (
	"encoding/json"
	"strings"

	"github.com/kyoh86/gogh/v2"
)

type ProjectFormat interface {
	Format(p gogh.Project) (string, error)
}

type ProjectFormatFunc func(gogh.Project) (string, error)

func (f ProjectFormatFunc) Format(p gogh.Project) (string, error) {
	return f(p)
}

var ProjectFormatFullFilePath = ProjectFormatFunc(func(p gogh.Project) (string, error) {
	return p.FullFilePath(), nil
})

var ProjectFormatRelPath = ProjectFormatFunc(func(p gogh.Project) (string, error) {
	return p.RelPath(), nil
})

var ProjectFormatRelFilePath = ProjectFormatFunc(func(p gogh.Project) (string, error) {
	return p.RelFilePath(), nil
})

var ProjectFormatURL = ProjectFormatFunc(func(p gogh.Project) (string, error) {
	return p.URL(), nil
})

var ProjectFormatJSON = ProjectFormatFunc(func(p gogh.Project) (string, error) {
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
})

func ProjectFormatFields(s string) ProjectFormat {
	return ProjectFormatFunc(func(p gogh.Project) (string, error) {
		return strings.Join([]string{
			p.FullFilePath(),
			p.RelFilePath(),
			p.URL(),
			p.RelPath(),
			p.Host(),
			p.Owner(),
			p.Name(),
		}, s), nil
	})
}
