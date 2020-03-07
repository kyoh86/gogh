// Code generated by main.go DO NOT EDIT.

package config

func Merge(envar Envar, cache CacheFile, keyring Keyring, config ConfigFile) (merged Merged) {
	merged.githubToken = new(GithubToken).Default().(string)
	merged.githubHost = new(GithubHost).Default().(string)
	merged.githubUser = new(GithubUser).Default().(string)
	merged.roots = new(Roots).Default().([]string)

	if config.GithubHost.Value().(string) != "" {
		merged.githubHost = config.GithubHost.Value().(string)
	}
	if len(config.Roots.Value().([]string)) != 0 {
		merged.roots = config.Roots.Value().([]string)
	}

	if keyring.GithubToken.Value().(string) != "" {
		merged.githubToken = keyring.GithubToken.Value().(string)
	}

	if cache.GithubUser.Value().(string) != "" {
		merged.githubUser = cache.GithubUser.Value().(string)
	}

	if envar.GithubToken.Value().(string) != "" {
		merged.githubToken = envar.GithubToken.Value().(string)
	}
	if envar.GithubHost.Value().(string) != "" {
		merged.githubHost = envar.GithubHost.Value().(string)
	}
	if len(envar.Roots.Value().([]string)) != 0 {
		merged.roots = envar.Roots.Value().([]string)
	}

	return
}

type Merged struct {
	githubToken string
	githubHost  string
	githubUser  string
	roots       []string
}

func (m *Merged) GithubToken() string {
	return m.githubToken
}

func (m *Merged) GithubHost() string {
	return m.githubHost
}

func (m *Merged) GithubUser() string {
	return m.githubUser
}

func (m *Merged) Roots() []string {
	return m.roots
}