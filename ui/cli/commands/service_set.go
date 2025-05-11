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
	defaultNameService  repository.DefaultNameService
	referenceParser     repository.ReferenceParser
	tokenService        auth.TokenService
	hostingService      hosting.HostingService
	finderService       workspace.FinderService
	workspaceService    workspace.WorkspaceService
	authenticateService auth.AuthenticateService
	flags               *config.Flags
	gitService          gitcore.GitService
}

func NewServiceSet(
	defaultNameService repository.DefaultNameService,
	tokenService auth.TokenService,
	workspaceService workspace.WorkspaceService,
	flags *config.Flags,
) *ServiceSet {
	return &ServiceSet{
		defaultNameService:  defaultNameService,
		referenceParser:     repository.NewReferenceParser(defaultNameService.GetDefaultHostAndOwner()),
		tokenService:        tokenService,
		hostingService:      github.NewHostingService(tokenService),
		finderService:       filesystem.NewFinderService(),
		workspaceService:    workspaceService,
		authenticateService: github.NewAuthenticateService(),
		flags:               flags,
		gitService:          gitimpl.NewService(),
	}
}
