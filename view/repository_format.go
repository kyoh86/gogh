package view

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v3"
)

type RepositoryPrinter interface {
	Print(p gogh.Repository) error
	Close() error
}

type RepositoryFuncPrinter func(gogh.Repository) error

func (f RepositoryFuncPrinter) Print(r gogh.Repository) error {
	return f(r)
}

func (f RepositoryFuncPrinter) Close() error {
	return nil
}

func NewRepositorySpecPrinter(w io.Writer) RepositoryPrinter {
	return RepositoryFuncPrinter(func(r gogh.Repository) error {
		fmt.Fprintln(w, r.Spec.String())
		return nil
	})
}

func NewRepositoryURLPrinter(w io.Writer) RepositoryPrinter {
	return RepositoryFuncPrinter(func(r gogh.Repository) error {
		fmt.Fprintln(w, r.URL)
		return nil
	})
}

func NewRepositoryJSONPrinter(w io.Writer) RepositoryPrinter {
	return RepositoryFuncPrinter(func(r gogh.Repository) error {
		buf, _ := json.Marshal(r)
		fmt.Fprintln(w, string(buf))
		return nil
	})
}
