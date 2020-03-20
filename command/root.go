package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Roots prints a gogh.root
func Roots(ev gogh.Env, all bool) error {
	if !all {
		fmt.Println(gogh.PrimaryRoot(ev))
		return nil
	}
	log.Println("info: Finding all roots...")
	for _, root := range ev.Roots() {
		fmt.Println(root)
	}
	return nil
}
