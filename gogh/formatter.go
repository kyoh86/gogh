package gogh

import (
	"fmt"
	"io"
)

// ProjectListFormatter holds project list to print them.
type ProjectListFormatter interface {
	Add(*Project)
	PrintAll(io.Writer, string) error
}

// ProjectListFormat specifies how gogh prints a project.
type ProjectListFormat string

// ProjectListFormat choices.
const (
	ProjectListFormatShort    = ProjectListFormat("short")
	ProjectListFormatFullPath = ProjectListFormat("full")
	ProjectListFormatRelPath  = ProjectListFormat("relative")
)

func (f ProjectListFormat) String() string {
	return string(f)
}

// ProjectListFormats shows all of ProjectListFormat constants.
func ProjectListFormats() []string {
	return []string{
		ProjectListFormatShort.String(),
		ProjectListFormatFullPath.String(),
		ProjectListFormatRelPath.String(),
	}
}

// Formatter will get a formatter to print list.
func (f ProjectListFormat) Formatter() (ProjectListFormatter, error) {
	switch f {
	case ProjectListFormatRelPath:
		return RelPathFormatter(), nil
	case ProjectListFormatFullPath:
		return FullPathFormatter(), nil
	case ProjectListFormatShort:
		return ShortFormatter(), nil
	}
	return nil, fmt.Errorf("%q is invalid project format", f)
}

// ShortFormatter prints each project as short as possible.
func ShortFormatter() ProjectListFormatter {
	return &shortListFormatter{
		dups: map[string]bool{},
	}
}

// FullPathFormatter prints each full-path of the project.
func FullPathFormatter() ProjectListFormatter {
	return &fullPathFormatter{&simpleCollector{}}
}

// RelPathFormatter prints each relative-path of the project
func RelPathFormatter() ProjectListFormatter {
	return &relPathFormatter{&simpleCollector{}}
}

type shortListFormatter struct {
	// mark duplicated subpath
	dups map[string]bool
	list []*Project
}

func (f *shortListFormatter) Add(r *Project) {
	for _, p := range r.Subpaths() {
		// (false, not ok) -> (false, ok) -> (true, ok) -> (true, ok) and so on
		_, f.dups[p] = f.dups[p]
	}
	f.list = append(f.list, r)
}

func (f *shortListFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.list {
		if _, err := fmt.Fprint(w, f.shortName(project)+sep); err != nil {
			return err
		}
	}
	return nil
}

func (f *shortListFormatter) shortName(r *Project) string {
	for _, p := range r.Subpaths() {
		if f.dups[p] {
			continue
		}
		return p
	}
	return r.FullPath
}

type simpleCollector struct {
	list []*Project
}

func (f *simpleCollector) Add(r *Project) {
	f.list = append(f.list, r)
}

type fullPathFormatter struct {
	*simpleCollector
}

func (f *fullPathFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.list {
		if _, err := fmt.Fprint(w, project.FullPath+sep); err != nil {
			return err
		}
	}
	return nil
}

type relPathFormatter struct {
	*simpleCollector
}

func (f *relPathFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.list {
		if _, err := fmt.Fprint(w, project.RelPath+sep); err != nil {
			return err
		}
	}
	return nil
}
