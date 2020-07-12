package command

import (
	"fmt"
	"os"
	"sync"

	"github.com/koron-go/prefixw"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/customio"
	"github.com/kyoh86/gogh/internal/git"
)

// Statuses shows all statuses of the gogh projects.
func Statuses(ev gogh.Env, gitClient GitClient, formatter gogh.ProjectListFormatter, primary bool, query string, detail bool) (retErr error) {
	walk := gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	if err := gogh.Query(ev, query, walk, func(p *gogh.Project) error {
		formatter.Add(p)
		return nil
	}); err != nil {
		return err
	}

	if err := formatter.Walk(func(p *gogh.Project, formatted string) error {
		if detail {
			return detailStatus(gitClient, p, formatted)
		}
		return summaryStatus(gitClient, p, formatted)
	}); err != nil {
		return err
	}
	return nil
}

func detailStatus(gitClient GitClient, p *gogh.Project, label string) error {
	var once sync.Once
	stdout := prefixw.New(&customio.LabeledWriter{
		Label: label,
		Once:  &once,
		Base:  os.Stdout,
	}, "\t")
	stderr := prefixw.New(&customio.LabeledWriter{
		Label: label,
		Once:  &once,
		Base:  os.Stderr,
	}, "\t")

	return gitClient.Status(p.FullPath, stdout, stderr)
}

func summaryStatus(gitClient GitClient, p *gogh.Project, label string) error {
	summary, err := gitClient.GetStatusSummary(p.FullPath, os.Stderr)
	if err != nil {
		return err
	}
	if summary != git.StatusSummaryClear {
		fmt.Printf("%s %s\n", summary, label)
	}
	return nil
}
