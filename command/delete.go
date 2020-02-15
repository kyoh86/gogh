package command

import (
	"fmt"
	"os"

	"github.com/kyoh86/ask"
	"github.com/kyoh86/gogh/gogh"
)

// Delete local projects
func Delete(ctx gogh.Context, primary bool, query string) error {
	var walk gogh.Walker = gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	var projects []*gogh.Project
	if err := gogh.Query(ctx, query, walk, func(p *gogh.Project) error {
		projects = append(projects, p)
		return nil
	}); err != nil {
		return err
	}

	if len(projects) == 0 {
		return fmt.Errorf("any projects did not matched for %q", query)
	}

	fmt.Println("Deleting projects. Please confirm them and answer by [y/n]")
	for _, p := range projects {
		fmt.Print(p.FullPath)
		yes, err := ask.Limit(5).YesNo()
		if err != nil {
			return err
		}
		if *yes {
			if err := os.RemoveAll(p.FullPath); err != nil {
				return err
			}
		}
	}
	return nil
}
