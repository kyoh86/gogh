package service

import (
	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/core/auth"
	gitcore "github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/filesystem"
	gitimpl "github.com/kyoh86/gogh/v3/infra/git"
	"github.com/kyoh86/gogh/v3/infra/github"
)

type ServiceSet struct {
	DefaultNameSource   string
	DefaultNameService  repository.DefaultNameService
	ReferenceParser     repository.ReferenceParser
	TokenSource         string
	TokenService        auth.TokenService
	HostingService      hosting.HostingService
	FinderService       workspace.FinderService
	WorkspaceSource     string
	WorkspaceService    workspace.WorkspaceService
	AuthenticateService auth.AuthenticateService
	FlagsSource         string
	Flags               *config.Flags
	GitService          gitcore.GitService
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
		DefaultNameSource:   defaultNameSource,
		DefaultNameService:  defaultNameService,
		ReferenceParser:     repository.NewReferenceParser(defaultNameService.GetDefaultHostAndOwner()),
		TokenSource:         tokenSource,
		TokenService:        tokenService,
		HostingService:      github.NewHostingService(tokenService),
		FinderService:       filesystem.NewFinderService(),
		WorkspaceSource:     workspaceSource,
		WorkspaceService:    workspaceService,
		AuthenticateService: github.NewAuthenticateService(),
		FlagsSource:         flagsSource,
		Flags:               flags,
		GitService:          gitimpl.NewService(),
	}
}
