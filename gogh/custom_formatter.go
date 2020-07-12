package gogh

import (
	"fmt"
	"io"
	"text/template"
)

type customListFormatter struct {
	template *template.Template
	*shortFormatter
}

func CustomFormatter(format string) (ProjectListFormatter, error) {
	// tmp.Funcs(
	short := &shortFormatter{
		dups:          map[string]bool{},
		coreFormatter: &coreFormatter{},
	}
	// NOTE: FuncMap is the type of the map defining the mapping from names to functions.
	// Each function must have either a single return value, or two return values of
	// which the second has type error. In that case, if the second (error)
	// return value evaluates to non-nil during execution, execution terminates and
	// Execute returns that error.
	//
	// When template execution invokes a function with an argument list, that list
	// must be assignable to the function's parameter types. Functions meant to
	// apply to arguments of arbitrary type can use parameters of type interface{} or
	// of type reflect.Value. Similarly, functions meant to return a result of arbitrary
	// type can return interface{} or reflect.Value.
	tmp, err := template.New("").Funcs(template.FuncMap{
		"short":    short.format,
		"full":     fullPath,
		"url":      fullURL,
		"relative": relPath,
		"null":     func() string { return "\x00" },
	}).Parse(format)
	if err != nil {
		return nil, fmt.Errorf("invalid custom format %w", err)
	}

	return &customListFormatter{
		template:       tmp,
		shortFormatter: short,
	}, nil
}

func (f *customListFormatter) format(w io.Writer, project *Project) error {
	return f.template.Execute(w, project)
}

func (f *customListFormatter) Add(r *Project) {
	// add to short list formatter (it has special "Add" func)
	f.shortFormatter.Add(r)
}

func (f *customListFormatter) Len() int {
	return f.shortFormatter.Len()
}

func (f *customListFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.shortFormatter.list {
		if err := f.format(w, project); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}
	}
	return nil
}
