package gogh

import (
	"fmt"
	"io"
)

// RepoListFormatter holds repository list to print them.
type RepoListFormatter interface {
	Add(*Local)
	PrintAll(io.Writer, string) error
}

// RepoListFormat specifies how gogh prints repo.
type RepoListFormat string

// RepoListFormat choices.
const (
	RepoListFormatShort    = RepoListFormat("short")
	RepoListFormatFullPath = RepoListFormat("full")
	RepoListFormatRelPath  = RepoListFormat("relative")
)

func (f RepoListFormat) String() string {
	return string(f)
}

// RepoListFormats shows all of RepoListFormat constants.
func RepoListFormats() []string {
	return []string{
		RepoListFormatShort.String(),
		RepoListFormatFullPath.String(),
		RepoListFormatRelPath.String(),
	}
}

// Formatter will get a formatter to print list.
func (f RepoListFormat) Formatter() (RepoListFormatter, error) {
	switch f {
	case RepoListFormatRelPath:
		return RelPathFormatter(), nil
	case RepoListFormatFullPath:
		return FullPathFormatter(), nil
	case RepoListFormatShort:
		return ShortFormatter(), nil
	}
	return nil, fmt.Errorf("%q is invalid repo format", f)
}

// ShortFormatter prints each repository as short as possible.
func ShortFormatter() RepoListFormatter {
	return &shortListFormatter{
		dups: map[string]bool{},
	}
}

// FullPathFormatter prints each full-path of the repository
func FullPathFormatter() RepoListFormatter {
	return &fullPathFormatter{&simpleCollector{}}
}

// RelPathFormatter prints each relative-path of the repository
func RelPathFormatter() RepoListFormatter {
	return &relPathFormatter{&simpleCollector{}}
}

type shortListFormatter struct {
	// mark duplicated subpath
	dups map[string]bool
	list []*Local
}

func (f *shortListFormatter) Add(r *Local) {
	for _, p := range r.Subpaths() {
		// (false, not ok) -> (false, ok) -> (true, ok) -> (true, ok) and so on
		_, f.dups[p] = f.dups[p]
	}
	f.list = append(f.list, r)
}

func (f *shortListFormatter) PrintAll(w io.Writer, sep string) error {
	for _, repo := range f.list {
		if _, err := fmt.Fprint(w, f.shortName(repo)+sep); err != nil {
			return err
		}
	}
	return nil
}

func (f *shortListFormatter) shortName(r *Local) string {
	for _, p := range r.Subpaths() {
		if f.dups[p] {
			continue
		}
		return p
	}
	return r.FullPath
}

type simpleCollector struct {
	list []*Local
}

func (f *simpleCollector) Add(r *Local) {
	f.list = append(f.list, r)
}

type fullPathFormatter struct {
	*simpleCollector
}

func (f *fullPathFormatter) PrintAll(w io.Writer, sep string) error {
	for _, repo := range f.list {
		if _, err := fmt.Fprint(w, repo.FullPath+sep); err != nil {
			return err
		}
	}
	return nil
}

type relPathFormatter struct {
	*simpleCollector
}

func (f *relPathFormatter) PrintAll(w io.Writer, sep string) error {
	for _, repo := range f.list {
		if _, err := fmt.Fprint(w, repo.RelPath+sep); err != nil {
			return err
		}
	}
	return nil
}
