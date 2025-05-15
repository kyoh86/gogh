package flags

import "github.com/spf13/cobra"

// BoolVarP binds a bool flag to a variable with a shorthand option.
func BoolVarP(cmd *cobra.Command, p *bool, name, shorthand string, value bool, usage string) {
	cmd.Flags().BoolVarP(p, name, shorthand, value, usage)
	cmd.Flags().Lookup(name).NoOptDefVal = "false"
}
