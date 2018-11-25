package gogh

import (
	"context"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) {
	run := func(label, conf string, f func(*testing.T)) {
		gitArgsForTest = []string{"--file", "-"}
		gitStdinForTest = []byte(conf)
		t.Run(label, f)
	}

	resetEnv := func(t *testing.T) {
		t.Helper()
		for _, key := range envNames {
			require.NoError(t, os.Setenv(key, ""))
		}
	}
	run("get current context", `[gogh "ghe"]
	host = foo.example.com
	host = bar.example.com
[gogh]
	log = Fatal
	user = kyoh86
	root = /go/src
	root = /foo/bar`, func(t *testing.T) {
		resetEnv(t)
		baseContext := context.WithValue(context.Background(), "test", "foo:bar")
		gotContext, err := CurrentContext(baseContext)
		require.NoError(t, err)
		assert.Equal(t, "foo:bar", gotContext.Value("test"))
		assert.Equal(t, "Fatal", gotContext.LogLevel())
		assert.Equal(t, "kyoh86", gotContext.UserName())
		assert.Equal(t, []string{"/go/src", "/foo/bar"}, gotContext.Roots())
		assert.Equal(t, "/go/src", gotContext.PrimaryRoot())
		assert.Equal(t, []string{"foo.example.com", "bar.example.com"}, gotContext.GHEHosts())
	})
	run("get log level from envar", `[gogh] log=dummy`, func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envLogLevel, "Warn"))
		log, err := getLogLevel()
		require.NoError(t, err)
		assert.Equal(t, "Warn", log)
	})

	run("get log level from git config", `[gogh] log=Error`, func(t *testing.T) {
		resetEnv(t)
		log, err := getLogLevel()
		require.NoError(t, err)
		assert.Equal(t, "Error", log)
	})

	run("get default log level from envar", ``, func(t *testing.T) {
		resetEnv(t)
		log, err := getLogLevel()
		require.NoError(t, err)
		assert.Equal(t, "Info", log)
	})

	run("expect to fail to get log level with invalid config", `[gogh] =foobar`, func(t *testing.T) {
		resetEnv(t)
		_, err := getLogLevel()
		assert.NotNil(t, err)
	})

	run("get user name from git config", "[gogh]\nuser = kyoh86", func(t *testing.T) {
		resetEnv(t)
		userName, err := getUserName()
		assert.NoError(t, err)
		assert.Equal(t, "kyoh86", userName)
	})

	run("get Github user name from envar", "", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envGithubUser, "kyoh87"))
		userName, err := getUserName()
		assert.NoError(t, err)
		assert.Equal(t, "kyoh87", userName)
	})

	run("get OS user name from envar", "", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envUserName, "kyoh88"))
		userName, err := getUserName()
		assert.NoError(t, err)
		assert.Equal(t, "kyoh88", userName)
	})

	run("expect to fail to get user name from anywhere", "", func(t *testing.T) {
		resetEnv(t)
		_, err := getUserName()
		assert.Error(t, err, "set gogh.user to your gitconfig")
	})

	run("expect to fail to get user with name invalid config", `[gogh] =foobar`, func(t *testing.T) {
		resetEnv(t)
		_, err := getUserName()
		assert.NotNil(t, err)
	})

	run("get root paths from envar", "[gogh]\nroot=/dummy", func(t *testing.T) {
		resetEnv(t)
		require.NoError(t, os.Setenv(envRoot, "/foo:/bar"))
		rts, err := getRoots()
		assert.NoError(t, err)
		assert.Equal(t, []string{"/foo", "/bar"}, rts)
	})

	run("get root paths from git config", "[gogh]\nroot=/foo\nroot=/bar", func(t *testing.T) {
		resetEnv(t)
		rts, err := getRoots()
		assert.NoError(t, err)
		assert.Equal(t, []string{"/foo", "/bar"}, rts)
	})

	run("get root paths from GOPATH", "", func(t *testing.T) {
		resetEnv(t)
		rts, err := getRoots()
		assert.NoError(t, err)
		assert.Equal(t, []string{filepath.Join(build.Default.GOPATH, "src")}, rts)
	})

	run("expect to fail to get roots with invalid config", `[gogh] =foobar`, func(t *testing.T) {
		resetEnv(t)
		_, err := getRoots()
		assert.NotNil(t, err)
	})

	run("get GHE hosts from git config", `[gogh "ghe"]`+"\nhost=foo.example.com\nhost=bar.example.com", func(t *testing.T) {
		resetEnv(t)
		hosts, err := getGHEHosts()
		assert.NoError(t, err)
		assert.Equal(t, []string{"foo.example.com", "bar.example.com"}, hosts)
	})

	run("expect to fail to get GHE hosts with invalid config", `[gogh] =foobar`, func(t *testing.T) {
		resetEnv(t)
		_, err := getGHEHosts()
		assert.NotNil(t, err, "expect to fail to get GHE hosts")
	})
}
