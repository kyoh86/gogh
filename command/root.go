package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Roots prints a gogh.root
func Roots(env gogh.Env, all bool) error {
	if !all {
		fmt.Println(gogh.PrimaryRoot(env))
		return nil
	}
	log.Println("info: Finding all roots...")
	for _, root := range env.Roots() {
		fmt.Println(root)
	}
	return nil
}
