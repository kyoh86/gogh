package flags

import "gopkg.in/alecthomas/kingpin.v2"

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
