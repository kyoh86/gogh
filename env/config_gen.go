package env

import (
	"fmt"
	types "github.com/kyoh86/appenv/types"
	"io"
)

// Code generated by main.go DO NOT EDIT.

type Config struct {
	yml YAML
}

func GetConfig(yamlReader io.Reader) (config Config, err error) {
	yml, err := loadYAML(yamlReader)
	if err != nil {
		return config, err
	}
	return buildConfig(yml)
}

func buildConfig(yml YAML) (config Config, err error) {
	config.yml = yml
	return
}

func (c *Config) Save(yamlWriter io.Writer) error {
	if err := saveYAML(yamlWriter, &c.yml); err != nil {
		return err
	}
	return nil
}

func PropertyNames() []string {
	return []string{"github.host", "github.user", "roots", "hooks"}
}

func (a *Config) Property(name string) (types.Config, error) {
	switch name {
	case "github.host":
		return &githubHostConfig{parent: a}, nil
	case "github.user":
		return &githubUserConfig{parent: a}, nil
	case "roots":
		return &rootsConfig{parent: a}, nil
	case "hooks":
		return &hooksConfig{parent: a}, nil
	}
	return nil, fmt.Errorf("invalid property name %q", name)
}

type githubHostConfig struct {
	parent *Config
}

func (a *githubHostConfig) Get() (string, error) {
	{
		p := a.parent.yml.GithubHost
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *githubHostConfig) Set(value string) error {
	{
		p := a.parent.yml.GithubHost
		if p == nil {
			p = new(GithubHost)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.GithubHost = p
	}
	return nil
}

func (a *githubHostConfig) Unset() {
	a.parent.yml.GithubHost = nil
}

type githubUserConfig struct {
	parent *Config
}

func (a *githubUserConfig) Get() (string, error) {
	{
		p := a.parent.yml.GithubUser
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *githubUserConfig) Set(value string) error {
	{
		p := a.parent.yml.GithubUser
		if p == nil {
			p = new(GithubUser)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.GithubUser = p
	}
	return nil
}

func (a *githubUserConfig) Unset() {
	a.parent.yml.GithubUser = nil
}

type rootsConfig struct {
	parent *Config
}

func (a *rootsConfig) Get() (string, error) {
	{
		p := a.parent.yml.Roots
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *rootsConfig) Set(value string) error {
	{
		p := a.parent.yml.Roots
		if p == nil {
			p = new(Roots)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Roots = p
	}
	return nil
}

func (a *rootsConfig) Unset() {
	a.parent.yml.Roots = nil
}

type hooksConfig struct {
	parent *Config
}

func (a *hooksConfig) Get() (string, error) {
	{
		p := a.parent.yml.Hooks
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *hooksConfig) Set(value string) error {
	{
		p := a.parent.yml.Hooks
		if p == nil {
			p = new(Hooks)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Hooks = p
	}
	return nil
}

func (a *hooksConfig) Unset() {
	a.parent.yml.Hooks = nil
}
