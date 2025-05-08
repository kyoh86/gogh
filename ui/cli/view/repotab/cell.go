package repotab

import (
	"strings"
	"time"

	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/morikuni/aec"
)

type CellBuilder interface {
	Build(hosting.Repository) (content string, style aec.ANSI)
}

type CellBuildFunc func(r hosting.Repository) (content string, style aec.ANSI)

func (f CellBuildFunc) Build(r hosting.Repository) (content string, style aec.ANSI) {
	return f(r)
}

var RepoRefCell = CellBuildFunc(func(r hosting.Repository) (content string, style aec.ANSI) {
	content = r.Ref.String()
	return content, aec.Bold
})

var DescriptionCell = CellBuildFunc(func(r hosting.Repository) (content string, style aec.ANSI) {
	content = r.Description
	return content, aec.DefaultF.With(aec.DefaultB)
})

var EmojiAttributesCell = CellBuildFunc(func(r hosting.Repository) (content string, style aec.ANSI) {
	var parts []string

	if r.Private {
		parts = append(parts, "🔒")
	}
	if r.Fork {
		parts = append(parts, "🔀")
	}
	if r.Archived {
		parts = append(parts, "🗃️")
	}

	return strings.Join(parts, " "), aec.EmptyBuilder.ANSI
})

var AttributesCell = CellBuildFunc(func(r hosting.Repository) (content string, style aec.ANSI) {
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

var UpdatedAtCell = CellBuildFunc(func(r hosting.Repository) (content string, style aec.ANSI) {
	return view.FuzzyAgoAbbr(time.Now(), r.UpdatedAt), aec.LightBlackF
})
