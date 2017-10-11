package cmd

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"github.com/kyoh86/ask"
	"github.com/kyoh86/gogh/internal/auth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// GitHubClient returns client for GitHub or error
func GitHubClient() (*github.Client, error) {
	ctx := context.Background()
	token := os.Getenv("GOGH_GITHUB_API_TOKEN")

	{
		client := auth.NewClient(ctx, token)
		if err := auth.Validate(ctx, client); err == nil {
			return client, nil
		}
	}

	// logrus.WithError(err).Debug("Failed to login GitHub with access token")
	username := os.Getenv("GOGH_GITHUB_USER")
	if username == "" {
		name, err := ask.Writer(os.Stderr).Message("GitHub login username").String()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to accept username")
		}
		username = *name
	}
	password := os.Getenv("GOGH_GITHUB_PASSWORD")
	if password == "" {
		pass, err := ask.Writer(os.Stderr).Message("GitHub login password").Hidden(true).String()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to accept password")
		}
		password = *pass
	}

	token, err := auth.Login(ctx, username, password)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to login GitHub with username/password")
	}
	logrus.WithField("token", token).Info("Created new API token. Set it in environment variable 'GOGH_GITHUB_API_TOKEN'")
	client := auth.NewClient(ctx, token)
	if err := auth.Validate(ctx, client); err != nil {
		return nil, errors.Wrap(err, "Failed to login GitHub with username/password")
	}
	return client, nil
}
