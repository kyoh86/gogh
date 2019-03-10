package gogh

var (
	envLogLevel        = "GOGH_LOG_LEVEL"
	envGoghGitHubUser  = "GOGH_GITHUB_USER"
	envGoghGitHubToken = "GOGH_GITHUB_TOKEN"
	envGoghGitHubHost  = "GOGH_GITHUB_HOST"
	envGitHubUser      = "GITHUB_USER"
	envGitHubToken     = "GITHUB_TOKEN"
	envGitHubHost      = "GITHUB_HOST"
	envGHEHosts        = "GOGH_GHE_HOST"
	envRoot            = "GOGH_ROOT"
	envNames           = []string{
		envLogLevel,
		envGoghGitHubUser,
		envGoghGitHubToken,
		envGoghGitHubHost,
		envGitHubUser,
		envGitHubToken,
		envGitHubHost,
		envUserName,
		envGHEHosts,
		envRoot,
	}
)
