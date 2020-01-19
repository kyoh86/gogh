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
	GetRemote(ctx gogh.Context, project *gogh.Project, name string) (*url.URL, error)
	GetRemotes(ctx gogh.Context, project *gogh.Project) (map[string]*url.URL, error)
	RenameRemote(ctx gogh.Context, project *gogh.Project, oldName, newName string) error
	RemoveRemote(ctx gogh.Context, project *gogh.Project, name string) error
	AddRemote(ctx gogh.Context, project *gogh.Project, name string, url *url.URL) error
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

func (i *mockGitClient) GetRemotes(ctx gogh.Context, project *gogh.Project) (map[string]*url.URL, error) {
	return nil, nil
}

func (i *mockGitClient) GetRemote(ctx gogh.Context, project *gogh.Project, name string) (*url.URL, error) {
	return nil, nil
}

func (i *mockGitClient) RenameRemote(ctx gogh.Context, project *gogh.Project, oldName, newName string) error {
	return nil
}

func (i *mockGitClient) RemoveRemote(ctx gogh.Context, project *gogh.Project, name string) error {
	return nil
}

func (i *mockGitClient) AddRemote(ctx gogh.Context, project *gogh.Project, name string, url *url.URL) error {
	return nil
}

var defaultGitClient gitClient = &internal.GitClient{}

func git() gitClient {
	return defaultGitClient
}
