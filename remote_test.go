package gogh_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	github "github.com/google/go-github/v33/github"
	testtarget "github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/mock_github"
	"github.com/wacul/ptr"
)

func TestRemoteController(t *testing.T) {
	ctx := context.Background()

	// TODO: list -> github.com/kyoh86/gogh, github.com/kyoh86/vim-gogh, github.kyoh86.dev/kyoh86/dotfiles
	// TODO: list --host github.com --owner kyoh86 -> github.com/kyoh86/gogh, github.com/kyoh86/vim-gogh

	// TODO: create gogh -> github.kyoh86.dev/kyoh86/gogh created (default: host=github.kyoh86.dev, user=kyoh86)
	// TODO: create kyoh86/gogh -> github.kyoh86.dev/kyoh86/gogh created (default: host=github.kyoh86.dev, user=kyoh86)
	// TODO: create github.com/kyoh86/gogh -> github.com/kyoh86/gogh created (default: host=github.kyoh86.dev, user=kyoh86)
	// NOTE: parsing gogh, kyoh86/gogh or github.com/kyoh86/gogh is the respoonsibilities of the "Descriptor" ->
	// NOTE: remove is the same for create
	t.Run("Unauthorized", func(t *testing.T) {
		remote := testtarget.NewRemoteController(testtarget.DefaultHost, "kyoh86")

		t.Run("List", func(t *testing.T) {
			t.Run("NilOption", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mock := mock_github.NewMockAdaptor(ctrl)
				remote.SetAdaptor(mock)

				mock.EXPECT().ListRepositories(ctx, "kyoh86", nil).Return([]*github.Repository{{
					Owner: &github.User{
						Login: ptr.String("kyoh86"),
					},
					Name: ptr.String("fake-1"),
				}, {
					Owner: &github.User{
						Login: ptr.String("kyoh86"),
					},
					Name: ptr.String("fake-2"),
				}}, nil, nil)
				projects, err := remote.List(ctx, nil)
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if len(projects) != 2 {
					t.Fatalf("expect some projects, but %d is gotten", len(projects))
				}
				expects := []string{
					"github.com/kyoh86/fake-1",
					"github.com/kyoh86/fake-2",
				}
				for i, expect := range expects {
					actual := projects[i].RelPath()
					if expect != actual {
						t.Errorf("expect project %q at %d but %q is gotten", expect, i, actual)
					}
				}
			})

			t.Run("EmptyOption", func(t *testing.T) {
				projects, err := remote.List(ctx, &testtarget.RemoteListOption{})
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if len(projects) <= 1 {
					t.Errorf("expect some projects, but %d is gotten", len(projects))
				}
			})

			t.Run("Organization", func(t *testing.T) {
				projects, err := remote.List(ctx, &testtarget.RemoteListOption{
					Organization: "kyoh86-tryouts",
				})
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if len(projects) != 2 {
					t.Fatalf("expect 2 projects, but %d is gotten", len(projects))
				}
				expect := map[string]struct{}{
					"github.com/kyoh86-tryouts/test": {},
					"github.com/kyoh86-tryouts/bare": {},
				}
				for _, p := range projects {
					_, match := expect[p.RelPath()]
					if !match {
						t.Errorf("unexpected project %q is gotten", p.RelPath())
					}
				}
				for p := range expect {
					t.Errorf("expected project %q, but not", p)
				}
			})

			t.Run("Query", func(t *testing.T) {
				projects, err := remote.List(ctx, &testtarget.RemoteListOption{
					Query: "kyoh86/gogh",
				})
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if len(projects) != 1 {
					t.Fatalf("expect one project, but %d is gotten", len(projects))
				}
				expect := "github.com/kyoh86/gogh"
				actual := projects[0].RelPath()
				if expect != actual {
					t.Errorf("expect project %q, but %q is gotten", expect, actual)
				}
			})

			t.Run("OrganizationAndQuery", func(t *testing.T) {
				projects, err := remote.List(ctx, &testtarget.RemoteListOption{
					Organization: "kyoh86-tryouts",
					Query:        "bare",
				})
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if len(projects) != 1 {
					t.Fatalf("expect one project, but %d is gotten", len(projects))
				}
				expect := "github.com/kyoh86-tryouts/bare"
				actual := projects[0].RelPath()
				if expect != actual {
					t.Errorf("expect project %q, but %q is gotten", expect, actual)
				}
			})

			t.Run("OrganizationAndQueryNoMatch", func(t *testing.T) {
				projects, err := remote.List(ctx, &testtarget.RemoteListOption{
					Organization: "kyoh86-tryouts",
					Query:        "no-match",
				})
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if len(projects) != 0 {
					t.Fatalf("expect zero project, but %d is gotten", len(projects))
				}
			})
		})
	})

	// func (r *RemoteController) Create(ctx context.Context, description Description) (Project, error)
	// func (r *RemoteController) Remove(ctx context.Context, description Description) error
}
