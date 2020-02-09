package command_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kyoh86/gogh/internal/context"
	"github.com/stretchr/testify/require"
)

type testService struct {
	gitCtrl   *gomock.Controller
	hubCtrl   *gomock.Controller
	gitClient *MockGitClient
	hubClient *MockHubClient
	root      string
	ctx       *context.MockContext
}

func (s testService) tearDown(t *testing.T) {
	t.Helper()
	s.gitCtrl.Finish()
	s.hubCtrl.Finish()
	require.NoError(t, os.RemoveAll(s.root))
}

func initTest(t *testing.T) *testService {
	t.Helper()
	gitCtrl := gomock.NewController(t)
	hubCtrl := gomock.NewController(t)
	gitClient := NewMockGitClient(gitCtrl)
	hubClient := NewMockHubClient(hubCtrl)

	root, err := ioutil.TempDir(os.TempDir(), "gogh-test")
	require.NoError(t, err)
	ctx := &context.MockContext{
		MRoot:       []string{root},
		MGitHubHost: "github.com",
	}

	return &testService{
		gitCtrl:   gitCtrl,
		hubCtrl:   hubCtrl,
		gitClient: gitClient,
		hubClient: hubClient,
		root:      root,
		ctx:       ctx,
	}
}
