package flags

import (
	"fmt"
	"regexp"
)

var (
	// BaseBranchRegexp for text intending a base-branch
	BaseBranchRegexp = regexp.MustCompile(`^(?:(?P<owner>` + NamePattern + `):)?(?P<branch>` + NamePattern + `)$`)
)

// BaseBranch contains identifier of pull-request target branch: "[owner:]branch"
type BaseBranch struct {
	Owner  string
	Branch string
}

// Set a value string to BaseBranch
func (r *BaseBranch) Set(value string) error {
	names := BaseBranchRegexp.SubexpNames()
	match := BaseBranchRegexp.FindStringSubmatch(value)
	if len(match) < len(names) {
		return fmt.Errorf("specified parameter '%s' is not a branch", value)
	}
	for i, name := range names {
		if match[i] == "" {
			continue
		}
		switch name {
		case "owner":
			r.Owner = match[i]
		case "branch":
			r.Branch = match[i]
		}
	}
	return nil
}

func (r *BaseBranch) String() string {
	if r.Owner == "" {
		return r.Branch
	}
	return fmt.Sprintf("%s:%s", r.Owner, r.Branch)
}
