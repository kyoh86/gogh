package command

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/command/internal"
	"github.com/kyoh86/gogh/gogh"
)

type gitClient interface {
	Init(ctx gogh.Context, project *gogh.Project, bare bool, template, separateGitDir string, shared gogh.ProjectShared) error
	Clone(ctx gogh.Context, project *gogh.Project, remote *url.URL, shallow bool) error
	Update(ctx gogh.Context, project *gogh.Project) error
}

type mockGitClient struct {
}

func (i *mockGitClient) Init(ctx gogh.Context, project *gogh.Project, bare bool, template, separateGitDir string, shared gogh.ProjectShared) error {
	return os.MkdirAll(filepath.Join(project.FullPath, ".git"), 0755)
}

func (i *mockGitClient) Clone(ctx gogh.Context, project *gogh.Project, remote *url.URL, shallow bool) error {
	return os.MkdirAll(filepath.Join(project.FullPath, ".git"), 0755)
}

func (i *mockGitClient) Update(ctx gogh.Context, project *gogh.Project) error {
	return nil
}

var defaultGitClient gitClient = &internal.GitClient{}

func git() gitClient {
	return defaultGitClient
}
