package gogh

import (
	"context"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Context holds configurations and environments
type Context interface {
	context.Context
	Stdout() io.Writer
	Stderr() io.Writer
	LogLevel() string

	GitHubUser() string
	GitHubToken() string
	GitHubHost() string
	Roots() []string
}

func PrimaryRoot(ctx Context) string {
	rts := ctx.Roots()
	return rts[0]
}

// Config holds configuration file values.
type Config map[string]interface{}

// CurrentContext get current context from OS envars and Git configurations
func CurrentContext(ctx context.Context, config Config) (Context, error) {
	gitHubUser := getGitHubUser(config)
	if gitHubUser == "" {
		// Make the error if the GitHubUser is not set.
		return nil, fmt.Errorf("failed to find user name. set %s in environment variable", envGoghGitHubUser)
	}
	if err := ValidateOwner(gitHubUser); err != nil {
		return nil, err
	}
	gitHubToken := getGitHubToken(config)
	gitHubHost := getGitHubHost(config)
	logLevel := getLogLevel(config)
	roots := getRoots(config)
	if err := validateRoots(roots); err != nil {
		return nil, err
	}
	return &implContext{
		Context:     ctx,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
		gitHubUser:  gitHubUser,
		gitHubToken: gitHubToken,
		gitHubHost:  gitHubHost,
		logLevel:    logLevel,
		roots:       roots,
	}, nil
}

type implContext struct {
	context.Context
	stdout      io.Writer
	stderr      io.Writer
	gitHubUser  string
	gitHubToken string
	gitHubHost  string
	logLevel    string
	roots       []string
}

func (c *implContext) Stdout() io.Writer {
	return c.stdout
}

func (c *implContext) Stderr() io.Writer {
	return c.stderr
}

func (c *implContext) GitHubUser() string {
	return c.gitHubUser
}

func (c *implContext) GitHubToken() string {
	return c.gitHubToken
}

func (c *implContext) GitHubHost() string {
	return c.gitHubHost
}

func (c *implContext) LogLevel() string {
	return c.logLevel
}

func (c *implContext) Roots() []string {
	return c.roots
}

func getConf(def string, envar string, config Config, key string, altEnvars ...string) string {
	if val := os.Getenv(envar); val != "" {
		log.Printf("debug: Context %s from envar is %q", key, val)
		return val
	}
	if config != nil {
		val, ok := config[key]
		if ok {
			str, ok := val.(string)
			if ok {
				log.Printf("debug: Context %s from file is %q", key, str)
				return str
			} else {
				log.Printf("warn: Config.%s expects string", confKeyRoot)
			}
		}
	}
	for _, e := range altEnvars {
		if val := os.Getenv(e); val != "" {
			log.Printf("debug: Context %s from alt-envar is %q", key, val)
			return val
		}
	}
	return def
}

func getGitHubToken(config Config) string {
	return getConf("", envGoghGitHubToken, config, confKeyGitHubToken, envGitHubToken)
}

func getGitHubHost(config Config) string {
	return getConf(DefaultHost, envGoghGitHubHost, config, confKeyGitHubHost, envGitHubHost)
}

func getGitHubUser(config Config) string {
	return getConf("", envGoghGitHubUser, config, confKeyGitHubUser, envGitHubUser, envUserName)
}

func getLogLevel(config Config) string {
	return getConf(DefaultLogLevel, envGoghLogLevel, config, confKeyLogLevel)
}

func getRoots(config Config) []string {
	envRoot := os.Getenv(envGoghRoot)
	if envRoot == "" {
		if config != nil {
			val, ok := config[confKeyRoot]
			if ok {
				arr, ok := val.([]string)
				if ok {
					return unique(arr)
				} else {
					log.Printf("warn: Config.%s expects string array", confKeyRoot)
				}
			}
		}
		gopaths := filepath.SplitList(build.Default.GOPATH)
		roots := make([]string, 0, len(gopaths))
		for _, gopath := range gopaths {
			roots = append(roots, filepath.Join(gopath, "src"))
		}
		return unique(roots)
	}
	return unique(filepath.SplitList(envRoot))
}

func validateRoots(roots []string) error {
	for i, v := range roots {
		path := filepath.Clean(v)
		_, err := os.Stat(path)
		switch {
		case err == nil:
			roots[i], err = filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
		case os.IsNotExist(err):
			roots[i] = path
		default:
			return err
		}
	}

	return nil
}
