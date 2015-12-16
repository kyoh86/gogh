package flags

import (
	"html/template"
	"time"

	"github.com/leekchan/gtf"
)

var formatter = gtf.New("gogh-flags-format").Funcs(template.FuncMap{
	"date": func(format string, t *time.Time) string {
		defer recover()

		if t == nil {
			return "-"
		}
		return t.In(time.Local).Format(format)
	},
})

// Template for a flag
func Template() *template.Template {
	return formatter
}
