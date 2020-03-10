package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Root prints a gogh.root
func Root(ctx gogh.Env, all bool) error {
	if !all {
		fmt.Println(gogh.PrimaryRoot(ctx))
		return nil
	}
	log.Println("info: Finding all roots...")
	for _, root := range ctx.Roots() {
		fmt.Println(root)
	}
	return nil
}
