package gh

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/tcnksm/go-gitconfig"
)

// Repository contains identifier of github repository: "owner/name"
type Repository struct {
	Owner string
	Repo  string
}

const (
	// NamePattern for owner/repo name
	NamePattern = `[a-zA-Z0-9\._-]+`
	// RepositoryPattern for text intending repository
	RepositoryPattern = `(?:(?:https?://)?github\.com/)?(?P<owner>` + NamePattern + `)/(?P<repo>` + NamePattern + `?)(?:\.git)?`
)

var (
	// RepositoryRegexp for text intending repository
	RepositoryRegexp = regexp.MustCompile(`^` + RepositoryPattern + `$`)
)

// Set a value string to Repository
func (r *Repository) Set(value string) error {
	names := RepositoryRegexp.SubexpNames()
	match := RepositoryRegexp.FindStringSubmatch(value)
	if len(match) < len(names) {
		return fmt.Errorf("specified parameter '%s' is not a repository", value)
	}
	for i, name := range names {
		if match[i] == "" {
			continue
		}
		switch name {
		case "owner":
			r.Owner = match[i]
		case "repo":
			r.Repo = match[i]
		}
	}
	return nil
}

func (r *Repository) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Repo)
}

var (
	workingOnce       sync.Once
	workingFound      bool
	workingRepository Repository
)

// WorkingRepository : Running on directory being a GitHub repository, get that owner/repo name
func WorkingRepository() (Repository, bool) {
	workingOnce.Do(func() {
		originURL, err := gitconfig.OriginURL()
		if err != nil {
			logrus.Debugf("gitconfig.OriginURL: %s", err.Error())
			return
		}

		if err := workingRepository.Set(originURL); err != nil {
			logrus.Debugf("ParseRepository(originURL): %s", err.Error())
			return
		}
		logrus.Debugf("CurrentOnRepository: %s", workingRepository.String())
		workingFound = true
	})

	return workingRepository, workingFound
}
