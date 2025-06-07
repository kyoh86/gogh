package overlay

import (
	"iter"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/typ"
)

func ForReference(ovs iter.Seq2[*Overlay, error], ref repository.Reference) iter.Seq2[*Overlay, error] {
	return typ.FilterE(ovs, func(ov *Overlay) (bool, error) {
		return ov.Match(ref)
	})
}

func ForPattern(ovs iter.Seq2[*Overlay, error], repoPattern string) iter.Seq2[*Overlay, error] {
	return typ.FilterE(ovs, func(ov *Overlay) (bool, error) {
		return ov.RepoPattern == repoPattern, nil
	})
}
