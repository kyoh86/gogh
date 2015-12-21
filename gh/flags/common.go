// Package flags defines common flags (name, description) of each commands
package flags

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/tcnksm/go-gitconfig"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	currentOnRepository bool
	currentOwner        string
	currentRepos        string
	currentOnce         sync.Once
)

// Running on directory being a GitHub repository, get that owner/repos name
func getCurrentRepositoryInfo() {
	currentOnce.Do(func() {
		currentOnRepository = false
		currentURL, err := gitconfig.OriginURL()
		if err != nil {
			logrus.Debugf("gitconfig.OriginURL: %s", err.Error())
			return
		}
		owner, repos, err := ParseRepository(currentURL)
		if err != nil {
			logrus.Debugf("ParseRepository(currentURL): %s", err.Error())
			return
		}
		currentOnRepository = true
		currentOwner = *owner
		currentRepos = *repos
		logrus.Debugf("CurrentOnRepository: %s/%s", *owner, *repos)
	})
}

const (
	// NamePattern for owner/repos name
	NamePattern = `[a-zA-Z0-9\._-]+`
)

var (
	// NameRegexp for owner/repos name
	NameRegexp = regexp.MustCompile(`^` + NamePattern + `$`)
	// RepositoryRegexp for text intending repository
	RepositoryRegexp = regexp.MustCompile(`^(?:(?:https?://)?github\.com/)?(?P<owner>` + NamePattern + `)/(?P<repos>` + NamePattern + `?)(?:\.git)?$`)
)

// ErrNotRepository returns error "specified parameter is not a repository"
func ErrNotRepository(text string) error {
	return fmt.Errorf("specified parameter '%s' is not a repository", text)
}

// ParseRepository : Parse a text intending repository, getting that owner and name
func ParseRepository(text string) (owner *string, repos *string, err error) {
	names := RepositoryRegexp.SubexpNames()
	match := RepositoryRegexp.FindStringSubmatch(text)
	if len(match) < len(names) {
		return nil, nil, ErrNotRepository(text)
	}
	for i, name := range names {
		switch name {
		case "owner":
			owner = &match[i]
		case "repos":
			repos = &match[i]
		}
	}
	return
}

// Owner sets flag for repository owner name
func Owner(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	getCurrentRepositoryInfo()
	f := cmd.Flag("owner", "Repository owner name").Short('o')
	if currentOnRepository {
		return f.Default(currentOwner)
	}
	return f
}

// Repos sets flag for repository name
func Repos(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	getCurrentRepositoryInfo()
	f := cmd.Flag("repos", "Repository name").Short('r')
	if currentOnRepository {
		return f.Default(currentRepos)
	}
	return f.Required()
}

// RepositoryValidator get validator for repository identifier
func RepositoryValidator(owner, repos *string) kingpin.Action {
	if owner == nil || repos == nil {
		panic("repository validator called with nil pointers")
	}
	return func(*kingpin.ParseContext) error {
		o, r, e := ParseRepository(*repos)
		if e == nil {
			*owner = *o
			*repos = *r
			return nil
		}
		if NameRegexp.MatchString(*owner) && NameRegexp.MatchString(*repos) {
			return nil
		}
		return ErrNotRepository(*owner + "/" + *repos)
	}
}

// Repository sets flags for specified pointer owner/repos
func Repository(cmd *kingpin.CmdClause, owner, repos *string) {
	Owner(cmd).StringVar(owner)
	Repos(cmd).StringVar(repos)
	cmd.PreAction(RepositoryValidator(owner, repos))
}

// Sort sets flag what to sort results by
func Sort(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	return cmd.Flag("sort", "What to sort results by")
}

// Direction sets flag the direction of the sort
func Direction(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	return cmd.Flag("direction", "The direction of the sort")
}

// PerPage sets flag specifies further pages
func PerPage(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	return cmd.Flag("per-page", "Specify further pages").Default("30")
}

// Page sets flag sets a custom page size up to 100
func Page(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	return cmd.Flag("page", "Custom page size up to 100").Default("1")
}
