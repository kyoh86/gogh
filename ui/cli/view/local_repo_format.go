package view

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v3/domain/local"
)

type LocalRepoFormat interface {
	Format(p local.Repo) (string, error)
}

type LocalRepoFormatFunc func(local.Repo) (string, error)

func (f LocalRepoFormatFunc) Format(p local.Repo) (string, error) {
	return f(p)
}

var LocalRepoFormatFullFilePath = LocalRepoFormatFunc(func(p local.Repo) (string, error) {
	return p.FullFilePath(), nil
})

var LocalRepoFormatRelPath = LocalRepoFormatFunc(func(p local.Repo) (string, error) {
	return p.RelPath(), nil
})

var LocalRepoFormatRelFilePath = LocalRepoFormatFunc(func(p local.Repo) (string, error) {
	return p.RelFilePath(), nil
})

var LocalRepoFormatURL = LocalRepoFormatFunc(func(p local.Repo) (string, error) {
	utxt, err := local.GetDefaultRemoteURL(context.Background(), p)
	if err != nil {
		if errors.Is(err, git.ErrRemoteNotFound) {
			utxt = "https://" + p.RelPath()
		} else {
			return "", err
		}
	}
	return utxt, nil
})

var LocalRepoFormatJSON = LocalRepoFormatFunc(func(p local.Repo) (string, error) {
	utxt, err := LocalRepoFormatURL(p)
	if err != nil {
		return "", err
	}
	buf, _ := json.Marshal(map[string]any{
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

func LocalRepoFormatFields(s string) LocalRepoFormat {
	return LocalRepoFormatFunc(func(p local.Repo) (string, error) {
		utxt, err := local.GetDefaultRemoteURL(context.Background(), p)
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
