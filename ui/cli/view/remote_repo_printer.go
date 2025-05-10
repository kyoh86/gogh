package view

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v3/core/hosting"
)

type RemoteRepoPrinter interface {
	Print(p hosting.Repository) error
	Close() error
}

type RemoteRepoFuncPrinter func(hosting.Repository) error

func (f RemoteRepoFuncPrinter) Print(r hosting.Repository) error {
	return f(r)
}

func (f RemoteRepoFuncPrinter) Close() error {
	return nil
}

func NewRemoteRepoRefPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r hosting.Repository) error {
		fmt.Fprintln(w, r.Ref.String())
		return nil
	})
}

func NewRemoteRepoURLPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r hosting.Repository) error {
		fmt.Fprintln(w, r.URL)
		return nil
	})
}

func NewRemoteRepoJSONPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r hosting.Repository) error {
		buf, _ := json.Marshal(r)
		fmt.Fprintln(w, string(buf))
		return nil
	})
}
