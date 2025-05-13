package view

import (
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v3/core/hosting"
)

type RepositoryPrinter interface {
	Print(p hosting.Repository) error
	Close() error
}

type RepositoryPrinterFunc func(hosting.Repository) error

func (f RepositoryPrinterFunc) Print(r hosting.Repository) error {
	return f(r)
}

func (f RepositoryPrinterFunc) Close() error {
	return nil
}

func FormatPrinter(w io.Writer, format hosting.RepositoryFormat) RepositoryPrinter {
	return RepositoryPrinterFunc(func(r hosting.Repository) error {
		s, err := format.Format(r)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, s)
		return nil
	})
}

func NewRepositoryPrinterRef(w io.Writer) RepositoryPrinter {
	return FormatPrinter(w, hosting.RepositoryFormatRef)
}

func NewRepositoryPrinterURL(w io.Writer) RepositoryPrinter {
	return FormatPrinter(w, hosting.RepositoryFormatURL)
}

func NewRepositoryPrinterJSON(w io.Writer) RepositoryPrinter {
	return FormatPrinter(w, hosting.RepositoryFormatJSON)
}
