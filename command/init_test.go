package command_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	code := m.Run()
	os.Exit(code)
}

type testService struct {
	gitCtrl   *gomock.Controller
	hubCtrl   *gomock.Controller
	envCtrl   *gomock.Controller
	gitClient *MockGitClient
	hubClient *MockHubClient
	root1     string
	root2     string
	env       *MockEnv
}

func (s testService) teardown(t *testing.T) {
	t.Helper()
	s.gitCtrl.Finish()
	s.hubCtrl.Finish()
	s.envCtrl.Finish()
	require.NoError(t, os.RemoveAll(s.root1))
	require.NoError(t, os.RemoveAll(s.root2))
}

func initTest(t *testing.T) *testService {
	t.Helper()
	gitCtrl := gomock.NewController(t)
	hubCtrl := gomock.NewController(t)
	envCtrl := gomock.NewController(t)
	gitClient := NewMockGitClient(gitCtrl)
	hubClient := NewMockHubClient(hubCtrl)
	ctxMock := NewMockEnv(envCtrl)

	root1, err := ioutil.TempDir(os.TempDir(), "gogh-test1")
	require.NoError(t, err)

	root2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
	require.NoError(t, err)

	ctxMock.EXPECT().GithubHost().AnyTimes().Return("github.com")
	ctxMock.EXPECT().Roots().AnyTimes().Return([]string{root1, root2})
	return &testService{
		gitCtrl:   gitCtrl,
		hubCtrl:   hubCtrl,
		envCtrl:   envCtrl,
		gitClient: gitClient,
		hubClient: hubClient,
		root1:     root1,
		root2:     root2,
		env:       ctxMock,
	}
}
