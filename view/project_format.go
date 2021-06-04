package view

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	git "github.com/go-git/go-git/v5"
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

func formatProjectURL(p gogh.Project) (string, error) {
	utxt, err := gogh.GetDefaultRemoteURLFromLocalProject(context.Background(), p)
	if err != nil {
		if errors.Is(err, git.ErrRemoteNotFound) {
			utxt = "https://" + p.RelPath()
		} else {
			return "", err
		}
	}
	return utxt, nil
}

var ProjectFormatURL = ProjectFormatFunc(func(p gogh.Project) (string, error) {
	return formatProjectURL(p)
})

var ProjectFormatJSON = ProjectFormatFunc(func(p gogh.Project) (string, error) {
	utxt, err := formatProjectURL(p)
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
