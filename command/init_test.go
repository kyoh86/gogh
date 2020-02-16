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
	gitClient *MockGitClient
	hubClient *MockHubClient
	root      string
	ctx       *MockContext
}

func (s testService) teardown(t *testing.T) {
	t.Helper()
	s.gitCtrl.Finish()
	s.hubCtrl.Finish()
	require.NoError(t, os.RemoveAll(s.root))
}

func initTest(t *testing.T) *testService {
	t.Helper()
	gitCtrl := gomock.NewController(t)
	hubCtrl := gomock.NewController(t)
	ctxCtrl := gomock.NewController(t)
	gitClient := NewMockGitClient(gitCtrl)
	hubClient := NewMockHubClient(hubCtrl)
	ctxMock := NewMockContext(ctxCtrl)

	root, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)

	return &testService{
		gitCtrl:   gitCtrl,
		hubCtrl:   hubCtrl,
		gitClient: gitClient,
		hubClient: hubClient,
		root:      root,
		ctx:       ctxMock,
	}
}
