package view

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v3"
)

type RemoteRepoPrinter interface {
	Print(p gogh.RemoteRepo) error
	Close() error
}

type RemoteRepoFuncPrinter func(gogh.RemoteRepo) error

func (f RemoteRepoFuncPrinter) Print(r gogh.RemoteRepo) error {
	return f(r)
}

func (f RemoteRepoFuncPrinter) Close() error {
	return nil
}

func NewRemoteRepoRefPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r gogh.RemoteRepo) error {
		fmt.Fprintln(w, r.Ref.String())
		return nil
	})
}

func NewRemoteRepoURLPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r gogh.RemoteRepo) error {
		fmt.Fprintln(w, r.URL)
		return nil
	})
}

func NewRemoteRepoJSONPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r gogh.RemoteRepo) error {
		buf, _ := json.Marshal(r)
		fmt.Fprintln(w, string(buf))
		return nil
	})
}
