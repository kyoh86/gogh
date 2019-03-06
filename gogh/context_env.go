package gogh

var (
	envLogLevel    = "GOGH_LOG_LEVEL"
	envGitHubUser  = "GITHUB_USER"
	envGitHubToken = "GITHUB_TOKEN"
	envRoot        = "GOGH_ROOT"
	envNames       = []string{
		envLogLevel,
		envGitHubUser,
		envGitHubToken,
		envUserName,
		envRoot,
	}
)
