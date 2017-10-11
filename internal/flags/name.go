package flags

import (
	"fmt"
	"regexp"
)

const (
	// NamePattern for text intending a name for GitHub
	NamePattern = `[^:\s'"` + "`" + `]+`
)

var (
	// NameRegexp for text intending a name for Github
	NameRegexp = regexp.MustCompile(`^` + NamePattern + `$`)
)

type Name string

// Set a value string to Name
func (n *Name) Set(value string) error {
	if !NameRegexp.MatchString(value) {
		return fmt.Errorf("specified parameter '%s' is not a name", value)
	}
	*n = Name(value)
	return nil
}

func (n *Name) String() string {
	return string(*n)
}
