package view

import (
	"context"
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
	return gogh.GetDefaultRemoteURLFromLocalProject(context.Background(), p)
})

var ProjectFormatJSON = ProjectFormatFunc(func(p gogh.Project) (string, error) {
	utxt, err := gogh.GetDefaultRemoteURLFromLocalProject(context.Background(), p)
	if err != nil {
		return "", err
	}
	buf, _ := json.Marshal(map[string]interface{}{
		"fullFilePath": p.FullFilePath(),
		"relFilePath":  p.RelFilePath(),
		"url":          utxt,
		"relPath":      p.RelPath(),
		"host":         p.Host(),
		"owner":        p.Owner(),
		"name":         p.Name(),
	})
	return string(buf), nil
})

func ProjectFormatFields(s string) ProjectFormat {
	return ProjectFormatFunc(func(p gogh.Project) (string, error) {
		utxt, err := gogh.GetDefaultRemoteURLFromLocalProject(context.Background(), p)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{
			p.FullFilePath(),
			p.RelFilePath(),
			utxt,
			p.RelPath(),
			p.Host(),
			p.Owner(),
			p.Name(),
		}, s), nil
	})
}
