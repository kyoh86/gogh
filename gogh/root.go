package gogh

import (
	"fmt"

	"github.com/kyoh86/gogh/repo"
)

// Root prints a gogh.root
func Root(all bool) error {
	if !all {
		fmt.Println(repo.PrimaryRoot())
		return nil
	}
	rts, err := repo.Roots()
	if err != nil {
		return err
	}
	for _, root := range rts {
		fmt.Println(root)
	}
	return nil
}
