package gogh

import (
	"fmt"
	"io"
)

// ProjectListFormatter holds project list to print them.
type ProjectListFormatter interface {
	Add(*Project)
	Len() int
	PrintAll(io.Writer, string) error
	format(io.Writer, *Project) error
}

// ShortFormatter prints each project as short as possible.
func ShortFormatter() ProjectListFormatter {
	return &shortListFormatter{
		dups:            map[string]bool{},
		simpleCollector: &simpleCollector{},
	}
}

// FullPathFormatter prints each full-path of the project.
func FullPathFormatter() ProjectListFormatter {
	return &fullPathFormatter{&simpleCollector{}}
}

// URLFormatter prints each project as url.
func URLFormatter() ProjectListFormatter {
	return &urlFormatter{&simpleCollector{}}
}

// RelPathFormatter prints each relative-path of the project
func RelPathFormatter() ProjectListFormatter {
	return &relPathFormatter{&simpleCollector{}}
}

type shortListFormatter struct {
	// mark duplicated subpath
	dups map[string]bool
	*simpleCollector
}

func (f *shortListFormatter) Add(r *Project) {
	for _, p := range r.Subpaths() {
		// (false, not ok) -> (false, ok) -> (true, ok) -> (true, ok) and so on
		_, f.dups[p] = f.dups[p]
	}
	f.simpleCollector.Add(r)
}

func (f *shortListFormatter) Len() int {
	return len(f.list)
}

func (f *shortListFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.list {
		if err := f.format(w, project); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}
	}
	return nil
}

func (f *shortListFormatter) format(w io.Writer, project *Project) error {
	_, err := fmt.Fprint(w, f.shortName(project))
	return err
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

func (f *simpleCollector) Len() int {
	return len(f.list)
}

type fullPathFormatter struct {
	*simpleCollector
}

func (f *fullPathFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.list {
		if err := f.format(w, project); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}
	}
	return nil
}

func (f *fullPathFormatter) format(w io.Writer, project *Project) error {
	_, err := fmt.Fprint(w, project.FullPath)
	return err
}

type urlFormatter struct {
	*simpleCollector
}

func (f *urlFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.list {
		if err := f.format(w, project); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}
	}
	return nil
}

func (f *urlFormatter) format(w io.Writer, project *Project) error {
	_, err := fmt.Fprint(w, "https://"+project.RelPath)
	return err
}

type relPathFormatter struct {
	*simpleCollector
}

func (f *relPathFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.list {
		if err := f.format(w, project); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}
	}
	return nil
}

func (f *relPathFormatter) format(w io.Writer, project *Project) error {
	_, err := fmt.Fprint(w, project.RelPath)
	return err
}
