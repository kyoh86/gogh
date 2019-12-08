package command

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/gogh"
)

// ProjectListFormat specifies how gogh prints a project.
type ProjectListFormat struct {
	label     string
	formatter gogh.ProjectListFormatter
}

func (f *ProjectListFormat) Set(value string) error {
	switch value {
	case ProjectListFormatLabelRelPath:
		f.formatter = gogh.RelPathFormatter()
		return nil
	case ProjectListFormatLabelFullPath:
		f.formatter = gogh.FullPathFormatter()
		return nil
	case ProjectListFormatLabelURL:
		f.formatter = gogh.URLFormatter()
		return nil
	case ProjectListFormatLabelShort:
		f.formatter = gogh.ShortFormatter()
		return nil
	}
	if strings.HasPrefix(value, "custom:") {
		er, err := gogh.CustomFormatter(strings.TrimPrefix(value, "custom:"))
		if err != nil {
			return fmt.Errorf("format custom: must have following valid template %w", err)
		}
		f.formatter = er
		return nil
	}
	return fmt.Errorf("format must be one of %s or 'custom:<advanced format>', got '%s'", strings.Join(ProjectListFormats(), ","), value)
}

func (f *ProjectListFormat) Formatter() gogh.ProjectListFormatter {
	return f.formatter
}

// ProjectListFormat choices.
const (
	ProjectListFormatLabelShort    = "short"
	ProjectListFormatLabelFullPath = "full"
	ProjectListFormatLabelURL      = "url"
	ProjectListFormatLabelRelPath  = "relative"
)

func (f ProjectListFormat) String() string {
	return f.label
}

// ProjectListFormats shows all of ProjectListFormat constants.
func ProjectListFormats() []string {
	return []string{
		ProjectListFormatLabelShort,
		ProjectListFormatLabelFullPath,
		ProjectListFormatLabelURL,
		ProjectListFormatLabelRelPath,
	}
}
