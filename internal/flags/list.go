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

// PageSize sets flag sets a custom page size up to 100
func PageSize(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	return cmd.Flag("per-page", "Custom page size up to 100").Default("30")
}

// Page sets flag specifies page number
func Page(cmd *kingpin.CmdClause) *kingpin.FlagClause {
	return cmd.Flag("page", "Specify page number").Default("1")
}
