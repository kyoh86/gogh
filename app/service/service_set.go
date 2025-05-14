package service

import (
	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/store"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

type ServiceSet struct {
	DefaultNameStore   store.Saver[repository.DefaultNameService]
	DefaultNameService repository.DefaultNameService

	TokenStore   store.Saver[auth.TokenService]
	TokenService auth.TokenService

	WorkspaceStore   store.Saver[workspace.WorkspaceService]
	WorkspaceService workspace.WorkspaceService

	FlagsStore store.Loader[*config.Flags]
	Flags      *config.Flags

	ReferenceParser     repository.ReferenceParser
	HostingService      hosting.HostingService
	FinderService       workspace.FinderService
	AuthenticateService auth.AuthenticateService
	GitService          git.GitService
}
