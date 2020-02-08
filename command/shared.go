package command

import (
	"fmt"
	"strconv"
)

// RepoShared notices file shared flag like group, all, everybody, 0766.
type RepoShared string

var validShared = map[string]struct{}{
	"false":     {},
	"true":      {},
	"umask":     {},
	"group":     {},
	"all":       {},
	"world":     {},
	"everybody": {},
}

// Set text as RepoShared
func (s *RepoShared) Set(text string) error {
	if _, ok := validShared[text]; ok {
		*s = RepoShared(text)
		return nil
	}
	if _, err := strconv.ParseInt(text, 8, 16); err == nil {
		*s = RepoShared(text)
		return nil
	}
	return fmt.Errorf(`invalid shared value %q; shared can be specified with "false", "true", "umask", "group", "all", "world", "everybody" or "0xxx" (octed value)`, text)
}

func (s RepoShared) String() string {
	return string(s)
}
