package view

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v3/domain/remote"
)

type RemoteRepoPrinter interface {
	Print(p remote.RemoteRepo) error
	Close() error
}

type RemoteRepoFuncPrinter func(remote.RemoteRepo) error

func (f RemoteRepoFuncPrinter) Print(r remote.RemoteRepo) error {
	return f(r)
}

func (f RemoteRepoFuncPrinter) Close() error {
	return nil
}

func NewRemoteRepoRefPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r remote.RemoteRepo) error {
		fmt.Fprintln(w, r.Ref.String())
		return nil
	})
}

func NewRemoteRepoURLPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r remote.RemoteRepo) error {
		fmt.Fprintln(w, r.URL)
		return nil
	})
}

func NewRemoteRepoJSONPrinter(w io.Writer) RemoteRepoPrinter {
	return RemoteRepoFuncPrinter(func(r remote.RemoteRepo) error {
		buf, _ := json.Marshal(r)
		fmt.Fprintln(w, string(buf))
		return nil
	})
}
