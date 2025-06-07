package service

import (
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/store"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

type ServiceSet struct {
	DefaultNameStore   store.Saver[repository.DefaultNameService]
	DefaultNameService repository.DefaultNameService

	TokenStore   store.Saver[auth.TokenService]
	TokenService auth.TokenService

	WorkspaceStore   store.Saver[workspace.WorkspaceService]
	WorkspaceService workspace.WorkspaceService

	FlagsStore store.Saver[*config.Flags]
	Flags      *config.Flags

	OverlayStore        overlay.OverlayStore
	OverlayService      overlay.OverlayService
	ReferenceParser     repository.ReferenceParser
	HostingService      hosting.HostingService
	FinderService       workspace.FinderService
	AuthenticateService auth.AuthenticateService
	GitService          git.GitService
}
