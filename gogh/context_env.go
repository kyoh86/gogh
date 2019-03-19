package gogh

var (
	envGoghLogLevel    = "GOGH_LOG_LEVEL"
	envGoghGitHubUser  = "GOGH_GITHUB_USER"
	envGoghGitHubToken = "GOGH_GITHUB_TOKEN"
	envGoghGitHubHost  = "GOGH_GITHUB_HOST"
	envGoghRoot        = "GOGH_ROOT"
	envGitHubUser      = "GITHUB_USER"
	envGitHubToken     = "GITHUB_TOKEN"
	envGitHubHost      = "GITHUB_HOST"
	envNames           = []string{
		envGoghLogLevel,
		envGoghGitHubUser,
		envGoghGitHubToken,
		envGoghGitHubHost,
		envGoghRoot,
		envGitHubUser,
		envGitHubToken,
		envGitHubHost,
		envUserName,
	}
)

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost     = "github.com"
	DefaultLogLevel = "warn"

	confKeyLogLevel    = "LogLevel"
	confKeyGitHubUser  = "GitHubUser"
	confKeyGitHubToken = "GitHubToken"
	confKeyGitHubHost  = "GitHubHost"
	confKeyRoot        = "Root"
)
