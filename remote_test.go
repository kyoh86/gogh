package gogh_test

import (
	"context"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestRemoteController(t *testing.T) {
	ctx := context.Background()

	t.Run("Unauthorized", func(t *testing.T) {
		remote := testtarget.NewRemoteController("github.com", "kyoh86")

		t.Run("List", func(t *testing.T) {
			t.Run("NilOption", func(t *testing.T) {
				projects, err := remote.List(ctx, nil)
				if err != nil {
					t.Fatalf("failed to listup: %s", err)
				}
				if len(projects) <= 1 {
					t.Errorf("expect some projects, but %d is gotten", len(projects))
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
