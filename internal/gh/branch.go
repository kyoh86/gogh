package gh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
)

var (
	// BranchRegexp for text intending a branch
	BranchRegexp = regexp.MustCompile(`^(?:` + RepositoryPattern + `:)?(?P<branch>` + NamePattern + `)$`)
)

// Branch contains identifier of github branch: "owner/name:branch"
type Branch struct {
	Repository
	Branch string
}

// Set a value string to Branch
func (r *Branch) Set(value string) error {
	names := BranchRegexp.SubexpNames()
	match := BranchRegexp.FindStringSubmatch(value)
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
		case "branch":
			r.Branch = match[i]
		}
	}
	return nil
}

func (r *Branch) String() string {
	return fmt.Sprintf("%s/%s:%s", r.Owner, r.Repo, r.Branch)
}

// WorkingBranch : Running on directory being a GitHub repository, get a current branch name of the working directory
func WorkingBranch() (Branch, bool) {
	repo, ok := WorkingRepository()
	if !ok {
		return Branch{}, false
	}

	branch := Branch{Repository: repo}
	var stdout bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				logrus.Debug("failed to get `git rev-parse for HEAD`")
				return branch, false
			}
		}
		logrus.Debugf("failed to get `git rev-parse for HEAD`: %s", exitError.Error())
		return branch, false
	}

	branch.Branch = strings.TrimRight(stdout.String(), "\n")
	return branch, true
}
