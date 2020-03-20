package gogh_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	code := m.Run()
	os.Exit(code)
}

type testService struct {
	envCtrl *gomock.Controller
	root1   string
	root2   string
	ev      *MockEnv
}

func initTest(t *testing.T) *testService {
	t.Helper()

	root1, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	root2, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)

	envCtrl := gomock.NewController(t)
	ev := NewMockEnv(envCtrl)
	ev.EXPECT().GithubUser().AnyTimes().Return("kyoh86")
	ev.EXPECT().GithubHost().AnyTimes().Return("github.com")
	ev.EXPECT().Roots().AnyTimes().Return([]string{root1, root2})

	return &testService{
		root1:   root1,
		root2:   root2,
		envCtrl: envCtrl,
		ev:      ev,
	}
}

func (s testService) teardown(t *testing.T) {
	require.NoError(t, os.RemoveAll(s.root1))
	require.NoError(t, os.RemoveAll(s.root2))
	s.envCtrl.Finish()
}
