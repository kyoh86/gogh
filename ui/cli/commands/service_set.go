package commands

import (
	"github.com/kyoh86/gogh/v3/core/auth"
	gitcore "github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/filesystem"
	gitimpl "github.com/kyoh86/gogh/v3/infra/git"
	"github.com/kyoh86/gogh/v3/infra/github"
)

type ServiceSet struct {
	defaultNameSource   string
	defaultNameService  repository.DefaultNameService
	referenceParser     repository.ReferenceParser
	tokenSource         string
	tokenService        auth.TokenService
	hostingService      hosting.HostingService
	finderService       workspace.FinderService
	workspaceSource     string
	workspaceService    workspace.WorkspaceService
	authenticateService auth.AuthenticateService
	flagsSource         string
	flags               *config.Flags
	gitService          gitcore.GitService
}

func NewServiceSet(
	defaultNameSource string,
	defaultNameService repository.DefaultNameService,
	tokenSource string,
	tokenService auth.TokenService,
	workspaceSource string,
	workspaceService workspace.WorkspaceService,
	flagsSource string,
	flags *config.Flags,
) *ServiceSet {
	return &ServiceSet{
		defaultNameSource:   defaultNameSource,
		defaultNameService:  defaultNameService,
		referenceParser:     repository.NewReferenceParser(defaultNameService.GetDefaultHostAndOwner()),
		tokenSource:         tokenSource,
		tokenService:        tokenService,
		hostingService:      github.NewHostingService(tokenService),
		finderService:       filesystem.NewFinderService(),
		workspaceSource:     workspaceSource,
		workspaceService:    workspaceService,
		authenticateService: github.NewAuthenticateService(),
		flagsSource:         flagsSource,
		flags:               flags,
		gitService:          gitimpl.NewService(),
	}
}
