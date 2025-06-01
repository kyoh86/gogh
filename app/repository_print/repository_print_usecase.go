package repository_print

import (
	"context"
	"fmt"
	"io"
	"iter"

	"github.com/kyoh86/gogh/v4/app/repository_print/repotab"
	"github.com/kyoh86/gogh/v4/core/hosting"
)

type UseCase struct {
	w      io.Writer
	format string
}

func repositoryFormatter(v string, w io.Writer) (RepositoryPrinter, error) {
	switch v {
	case "", "table":
		return repotab.NewPrinter(w, repotab.TermWidth(w), repotab.Styled(false, w)), nil
	case "ref":
		return NewRepositoryPrinterRef(w), nil
	case "url":
		return NewRepositoryPrinterURL(w), nil
	case "json":
		return NewRepositoryPrinterJSON(w), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}

func NewUseCase(w io.Writer, format string) *UseCase {
	return &UseCase{
		w:      w,
		format: format,
	}
}

func (uc *UseCase) Execute(_ context.Context, r iter.Seq2[*hosting.Repository, error]) error {
	printer, err := repositoryFormatter(uc.format, uc.w)
	if err != nil {
		return err
	}
	for repo, err := range r {
		if err != nil {
			return err
		}
		if repo == nil {
			continue
		}
		if err := printer.Print(*repo); err != nil {
			return err
		}
	}
	return printer.Close()
}

type RepositoryPrinterFunc func(hosting.Repository) error

func (f RepositoryPrinterFunc) Print(r hosting.Repository) error {
	return f(r)
}

func (f RepositoryPrinterFunc) Close() error {
	return nil
}

type RepositoryPrinter interface {
	Print(p hosting.Repository) error
	Close() error
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
