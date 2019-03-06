package gogh

import (
	"fmt"
	"strconv"
)

// ProjectShared notices file shared flag like group, all, everybody, 0766.
type ProjectShared string

var validShared = map[string]struct{}{
	"false":     {},
	"true":      {},
	"umask":     {},
	"group":     {},
	"all":       {},
	"world":     {},
	"everybody": {},
}

// Set text as ProjectShared
func (s *ProjectShared) Set(text string) error {
	if _, ok := validShared[text]; ok {
		*s = ProjectShared(text)
		return nil
	}
	if _, err := strconv.ParseInt(text, 8, 16); err == nil {
		*s = ProjectShared(text)
		return nil
	}
	return fmt.Errorf(`invalid shared value %q; shared can be specified with "false", "true", "umask", "group", "all", "world", "everybody" or "0xxx" (octed value)`, text)
}

func (s ProjectShared) String() string {
	return string(s)
}
