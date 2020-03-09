// Code generated by main.go DO NOT EDIT.

package env

import "io"

func GetAccess(yamlReader io.Reader, keyringService string, envarPrefix string) (access Access, err error) {
	yml, err := loadYAML(yamlReader)
	if err != nil {
		return access, err
	}
	keyring, err := loadKeyring(keyringService)
	if err != nil {
		return access, err
	}
	envar, err := getEnvar(envarPrefix)
	if err != nil {
		return access, err
	}
	access.roots = new(Roots).Default().([]string)
	if yml.Roots != nil {
		access.roots = yml.Roots.Value().([]string)
	}
	if envar.Roots != nil {
		access.roots = envar.Roots.Value().([]string)
	}

	access.githubHost = new(GithubHost).Default().(string)
	if yml.GithubHost != nil {
		access.githubHost = yml.GithubHost.Value().(string)
	}
	if envar.GithubHost != nil {
		access.githubHost = envar.GithubHost.Value().(string)
	}

	access.githubToken = new(GithubToken).Default().(string)
	if keyring.GithubToken != nil {
		access.githubToken = keyring.GithubToken.Value().(string)
	}
	if envar.GithubToken != nil {
		access.githubToken = envar.GithubToken.Value().(string)
	}

	return
}

type Access struct {
	roots       []string
	githubHost  string
	githubToken string
}

func (a *Access) Roots() []string {
	return a.roots
}

func (a *Access) GithubHost() string {
	return a.githubHost
}

func (a *Access) GithubToken() string {
	return a.githubToken
}
