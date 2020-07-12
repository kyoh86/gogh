package gogh

import (
	"fmt"
	"io"
	"sync"
)

// FormatWalkFunc is the function to accept projects
// and formatted string in the ProjectListFormatter.Walk
type FormatWalkFunc func(project *Project, formatted string) error

// ProjectListFormatter holds project list to print them.
type ProjectListFormatter interface {
	Add(*Project)
	Len() int
	PrintAll(io.Writer, string) error
	Walk(FormatWalkFunc) error
}

type coreFormatter struct {
	lock   sync.Mutex
	list   []*Project
	format func(*Project) string
}

func (f *coreFormatter) Add(project *Project) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.list = append(f.list, project)
}

func (f *coreFormatter) Len() int {
	return len(f.list)
}

func (f *coreFormatter) Walk(callback FormatWalkFunc) error {
	for _, project := range f.list {
		if err := callback(project, f.format(project)); err != nil {
			return err
		}
	}
	return nil
}

func (f *coreFormatter) PrintAll(w io.Writer, sep string) error {
	return f.Walk(func(project *Project, formatted string) error {
		if _, err := fmt.Fprint(w, formatted); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}
		return nil
	})
}

// ShortFormatter prints each project as short as possible.
func ShortFormatter() ProjectListFormatter {
	shorten := &shortFormatter{
		dups:          map[string]bool{},
		coreFormatter: &coreFormatter{},
	}

	shorten.coreFormatter.format = shorten.format
	return shorten
}

// FullPathFormatter prints each full-path of the project.
func FullPathFormatter() ProjectListFormatter {
	return &coreFormatter{format: fullPath}
}

func fullPath(project *Project) string {
	return project.FullPath
}

// URLFormatter prints each project as url.
func URLFormatter() ProjectListFormatter {
	return &coreFormatter{format: fullURL}
}

func fullURL(project *Project) string {
	return "https://" + project.RelPath
}

// RelPathFormatter prints each relative-path of the project
func RelPathFormatter() ProjectListFormatter {
	return &coreFormatter{format: relPath}
}

func relPath(project *Project) string {
	return project.RelPath
}

type shortFormatter struct {
	lock sync.Mutex
	// mark duplicated subpath
	dups map[string]bool
	*coreFormatter
}

func (f *shortFormatter) Add(project *Project) {
	f.lock.Lock()
	defer f.lock.Unlock()
	for _, p := range project.Subpaths() {
		// (false, not ok) -> (false, ok) -> (true, ok) -> (true, ok) and so on
		_, f.dups[p] = f.dups[p]
	}
	f.coreFormatter.Add(project)
}

func (f *shortFormatter) format(project *Project) string {
	for _, p := range project.Subpaths() {
		if f.dups[p] {
			continue
		}
		return p
	}
	return project.FullPath
}
