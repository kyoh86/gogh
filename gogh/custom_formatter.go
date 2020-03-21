package gogh

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type customListFormatter struct {
	template *template.Template
	*shortListFormatter
	*fullPathFormatter
	*urlFormatter
	*relPathFormatter
}

func CustomFormatter(format string) (ProjectListFormatter, error) {
	// share list
	simple := &simpleCollector{}

	// tmp.Funcs(
	short := &shortListFormatter{
		dups:            map[string]bool{},
		simpleCollector: simple,
	}
	full := &fullPathFormatter{
		simpleCollector: simple,
	}
	url := &urlFormatter{
		simpleCollector: simple,
	}
	rel := &relPathFormatter{
		simpleCollector: simple,
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
		"short":    formatToTemplateFunc(short.format),
		"full":     formatToTemplateFunc(full.format),
		"url":      formatToTemplateFunc(url.format),
		"relative": formatToTemplateFunc(rel.format),
		"null":     func() string { return "\x00" },
	}).Parse(format)
	if err != nil {
		return nil, fmt.Errorf("invalid custom format %w", err)
	}

	return &customListFormatter{
		template:           tmp,
		shortListFormatter: short,
		fullPathFormatter:  full,
		urlFormatter:       url,
		relPathFormatter:   rel,
	}, nil
}

func formatToTemplateFunc(format func(w io.Writer, project *Project) error) func(project *Project) (string, error) {
	return func(project *Project) (string, error) {
		buf := new(bytes.Buffer)
		if err := format(buf, project); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}
func (f *customListFormatter) format(w io.Writer, project *Project) error {
	return f.template.Execute(w, project)
}

func (f *customListFormatter) Add(r *Project) {
	// add to short list formatter (it has special "Add" func)
	f.shortListFormatter.Add(r)
}

func (f *customListFormatter) Len() int {
	return f.shortListFormatter.Len()
}

func (f *customListFormatter) PrintAll(w io.Writer, sep string) error {
	for _, project := range f.shortListFormatter.list {
		if err := f.format(w, project); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, sep); err != nil {
			return err
		}
	}
	return nil
}
