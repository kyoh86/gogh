package repotab

import (
	"strings"
	"time"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/view"
	"github.com/morikuni/aec"
)

type CellBuilder interface {
	Build(gogh.Repository) (content string, style aec.ANSI)
}

type CellBuildFunc func(r gogh.Repository) (content string, style aec.ANSI)

func (f CellBuildFunc) Build(r gogh.Repository) (content string, style aec.ANSI) {
	return f(r)
}

var SpecCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
	content = r.Spec.String()
	return content, aec.Bold
})

var DescriptionCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
	content = r.Description
	return content, aec.DefaultF.With(aec.DefaultB)
})

var EmojiAttributesCell = CellBuildFunc(func(r gogh.Repository) (content string, style aec.ANSI) {
	// UNDONE: this breaks terminal
	contents := []string{""}
	if r.Private {
		contents[0] = "\U0001F512\uFE0F "
	} else {
		contents[0] = ""
	}
	if r.Fork {
		contents = append(contents, "\U0001F500\uFE0F ")
	}
	if r.Archived {
		contents = append(contents, "\U0001F5C3\uFE0F ")
	}
	return strings.Join(contents, ""), aec.EmptyBuilder.ANSI
})

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
	if r.PushedAt == nil {
		return "", aec.LightBlackF
	}
	return view.FuzzyAgoAbbr(time.Now(), *r.PushedAt), aec.LightBlackF
})
