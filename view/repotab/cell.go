package repotab

import (
	"strings"
	"time"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/view"
	"github.com/morikuni/aec"
)

type Column struct {
	Truncatable bool
	MinWidth    int
	Elipsis     string

	Priority    int
	CellBuilder CellBuilder
}

type CellBuilder interface {
	Build(gogh.Repository) (content string, style aec.ANSI)
}

type CellBuildFunc func(r gogh.Repository) (content string, style aec.ANSI)

func (f CellBuildFunc) Build(r gogh.Repository) (content string, style aec.ANSI) {
	return f(r)
}

const GenericElipsis = ".."

var SpecCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
	content = r.Spec.String()
	return content, aec.Bold
})

var DescriptionCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
	content = r.Description
	return content, aec.DefaultF.With(aec.DefaultB)
})

// 	// UNDONE: this breaks aec.Apply
// var EmojiAttributesCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
// 	contents := []string{""}
// 	if r.Private {
// 		contents[0] = "🔒 "
// 	} else {
// 		contents[0] = ""
// 	}
// 	if r.Fork {
// 		contents = append(contents, "🔀 ")
// 	}
// 	if r.Archived {
// 		contents = append(contents, "🗃️ ")
// 	}
// 	return strings.Join(contents, ""), aec.EmptyBuilder.ANSI
// })

var AttributesCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
	contents := []string{""}
	if r.Private {
		style = aec.YellowF
		contents[0] = "private"
	} else {
		style = aec.LightBlackF
		contents[0] = "public"
	}
	if r.Fork {
		contents = append(contents, "fork")
	}
	if r.Archived {
		contents = append(contents, "archived")
	}
	return strings.Join(contents, ","), style
})

var PushedAtCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
	return view.FuzzyAgoAbbr(time.Now(), r.PushedAt), aec.LightBlackF
})
