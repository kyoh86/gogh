package cf

import (
	"github.com/kyoh86/gogh/conf"
	"github.com/kyoh86/gogh/gh"
	"gopkg.in/alecthomas/kingpin.v2"
)

// SetCommand will set a configuration
func SetCommand(c *kingpin.CmdClause) gh.Command {
	var (
		name  string
		value string
	)

	c.Flag("name", "Configuration name").Required().EnumVar(&name, conf.ConfigureItems...)
	c.Flag("value", "Configuration value").Required().StringVar(&value)

	return func() error {
		switch name {
		case conf.ConfigurationItemAccessToken:
			return conf.Set(func(c conf.Configures) (conf.Configures, error) {
				c.AccessToken = value
				return c, nil
			})
		default:
			return nil
		}
	}
}
