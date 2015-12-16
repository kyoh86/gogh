// Package flags defines common flags (name, description) of each commands
package flags

import "gopkg.in/alecthomas/kingpin.v2"

type cmd interface {
	Flag(name, help string) *kingpin.FlagClause
}

// Owner sets flag for repository owner name
func Owner(cmd cmd) *kingpin.FlagClause {
	return cmd.Flag("owner", "Repository owner name").Short('o')
}

// Repos sets flag for repository name
func Repos(cmd cmd) *kingpin.FlagClause {
	return cmd.Flag("repos", "Repository name").Short('r')
}

// Sort sets flag what to sort results by
func Sort(cmd cmd) *kingpin.FlagClause {
	return cmd.Flag("sort", "What to sort results by")
}

// Direction sets flag the direction of the sort
func Direction(cmd cmd) *kingpin.FlagClause {
	return cmd.Flag("direction", "The direction of the sort")
}

// PerPage sets flag specifies further pages
func PerPage(cmd cmd) *kingpin.FlagClause {
	return cmd.Flag("per-page", "Specify further pages").Default("30")
}

// Page sets flag sets a custom page size up to 100
func Page(cmd cmd) *kingpin.FlagClause {
	return cmd.Flag("page", "Custom page size up to 100").Default("1")
}
