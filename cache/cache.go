package cache

import (
	"io"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	yaml "gopkg.in/yaml.v2"
)

func GetCache(reader io.Reader) (cache Cache, err error) {
	yml, err := loadYAML(reader)
	if err != nil {
		return cache, err
	}
	cache.yml = yml
	return
}

var EmptyYAMLReader io.Reader = nil

func loadYAML(r io.Reader) (yml YAML, err error) {
	if r == EmptyYAMLReader {
		return
	}
	err = yaml.NewDecoder(r).Decode(&yml)
	return
}

type Cache struct {
	yml YAML
}

type YAML struct {
	GithubUser string `yaml:"githubUser"`
}

func (a *Cache) UnsetGithubUser(key, value string) {
	a.yml.GithubUser = ""
}

func (a *Cache) SetGithubUser(key, value string) {
	buf, err := jws.Sign([]byte(value), jwa.HS512, []byte(key))
	if err != nil {
		return
	}

	a.yml.GithubUser = string(buf)
}

func (a *Cache) GetGithubUser(key string) string {
	if a.yml.GithubUser == "" {
		return ""
	}

	verified, err := jws.Verify([]byte(a.yml.GithubUser), jwa.HS512, []byte(key))
	if err != nil {
		return ""
	}
	return string(verified)
}
