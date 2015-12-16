package cl

import (
	"github.com/google/go-github/github"
	"github.com/kyoh86/gogh/auth"
	"github.com/kyoh86/gogh/conf"
	"github.com/kyoh86/gogh/util"
)

// GitHubClient returns client for GitHub or error
func GitHubClient() (client *github.Client, err error) {
	return client, conf.Set(func(c conf.Configures) (conf.Configures, error) {
		client, err = auth.NewClient(c.AccessToken)
		if err == nil {
			return c, conf.ErrNotUpdated
		}

		// logrus.WithError(err).Debug("Failed to login GitHub with access token")
		username, err := Ask("GitHub login username")
		if err != nil {
			return c, util.WrapErr("Failed to accept username", err)
		}
		password, err := Secret("GitHub login password")
		if err != nil {
			return c, util.WrapErr("Failed to accept password", err)
		}

		token, err := auth.Login(username, password)
		if err != nil {
			return c, util.WrapErr("Failed to login GitHub with username/password", err)
		}
		client, err = auth.NewClient(token)
		if err != nil {
			return c, util.WrapErr("Failed to login GitHub with username/password", err)
		}

		c.AccessToken = token
		return c, nil
	})
}
