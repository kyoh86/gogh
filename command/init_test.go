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
	ctxCtrl   *gomock.Controller
	gitClient *MockGitClient
	hubClient *MockHubClient
	root1     string
	root2     string
	ctx       *MockContext
}

func (s testService) teardown(t *testing.T) {
	t.Helper()
	s.gitCtrl.Finish()
	s.hubCtrl.Finish()
	s.ctxCtrl.Finish()
	require.NoError(t, os.RemoveAll(s.root1))
	require.NoError(t, os.RemoveAll(s.root2))
}

func initTest(t *testing.T) *testService {
	t.Helper()
	gitCtrl := gomock.NewController(t)
	hubCtrl := gomock.NewController(t)
	ctxCtrl := gomock.NewController(t)
	gitClient := NewMockGitClient(gitCtrl)
	hubClient := NewMockHubClient(hubCtrl)
	ctxMock := NewMockContext(ctxCtrl)

	root1, err := ioutil.TempDir(os.TempDir(), "gogh-test1")
	require.NoError(t, err)

	root2, err := ioutil.TempDir(os.TempDir(), "gogh-test2")
	require.NoError(t, err)

	ctxMock.EXPECT().Done().AnyTimes().Return(nil)
	ctxMock.EXPECT().GitHubHost().AnyTimes().Return("github.com")
	ctxMock.EXPECT().GitHubUser().AnyTimes().Return("kyoh86")
	ctxMock.EXPECT().Root().AnyTimes().Return([]string{root1, root2})
	ctxMock.EXPECT().PrimaryRoot().AnyTimes().Return(root1)
	return &testService{
		gitCtrl:   gitCtrl,
		hubCtrl:   hubCtrl,
		ctxCtrl:   ctxCtrl,
		gitClient: gitClient,
		hubClient: hubClient,
		root1:     root1,
		root2:     root2,
		ctx:       ctxMock,
	}
}
