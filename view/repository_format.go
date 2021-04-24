package view

import (
	"encoding/json"
	"strings"

	"github.com/kyoh86/gogh/v2"
)

type RepositoryFormat interface {
	Format(p gogh.Repository) (string, error)
}

type RepositoryFormatFunc func(gogh.Repository) (string, error)

func (f RepositoryFormatFunc) Format(r gogh.Repository) (string, error) {
	return f(r)
}

func RepositoryFormatURL(r gogh.Repository) (string, error) {
	return r.URL, nil
}

func RepositoryFormatJSON(r gogh.Repository) (string, error) {
	buf, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func RepositoryFormatFields(s string) RepositoryFormat {
	return RepositoryFormatFunc(func(r gogh.Repository) (string, error) {
		return strings.Join([]string{
			r.Host(),
			r.Owner(),
			r.Name(),
		}, s), nil
	})
}
