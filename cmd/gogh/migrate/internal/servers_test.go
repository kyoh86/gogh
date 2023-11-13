package migrate_test

import (
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2/cmd/gogh/migrate/internal"
)

func TestServers(t *testing.T) {
	const (
		user1  = "kyoh86"
		user2  = "anonymous"
		host1  = "example.com" // host a not default
		host2  = "kyoh86.dev"  // host a not default
		token1 = "1111111111111111111111111111111111111111"
		token2 = "2222222222222222222222222222222222222222"
	)

	t.Run("Empty", func(t *testing.T) {
		var servers testtarget.Servers
		_, err := servers.Default()
		if !errors.Is(err, testtarget.ErrNoServer) {
			t.Errorf("expect error: %v, acutal: %v", testtarget.ErrNoServer, err)
		}

		list, err := servers.List()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list) != 0 {
			t.Fatalf("length mismatch: -want +got\n -%d\n +%d", 0, len(list))
		}
	})

	t.Run("Manipulate", func(t *testing.T) {
		var servers testtarget.Servers
		if err := servers.Set("", "", ""); err == nil {
			t.Fatalf("expect error, actual nil")
		}
		for _, testcase := range []struct {
			title string
			host  string
			user  string
			token string
		}{
			{
				title: "SetFirst",
				host:  host1,
				user:  user1,
				token: token1,
			},
			{
				title: "SetSecond",
				host:  host2,
				user:  user2,
				token: token2,
			},
		} {
			t.Run(testcase.title, func(t *testing.T) {
				if err := servers.Set(testcase.host, testcase.user, testcase.token); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				def, err := servers.Default()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if def.Host() != host1 {
					t.Errorf("expect host %q, actual: %q", host1, def.Host())
				}
				if def.User() != user1 {
					t.Errorf("expect user %q, actual: %q", user1, def.User())
				}
				if def.Token() != token1 {
					t.Errorf("expect token %q, actual: %q", token1, def.Token())
				}

				added, err := servers.Find(testcase.host)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if added.Host() != testcase.host {
					t.Errorf("expect host %q, actual: %q", testcase.host, added.Host())
				}
				if added.User() != testcase.user {
					t.Errorf("expect user %q, actual: %q", testcase.user, added.User())
				}
				if added.Token() != testcase.token {
					t.Errorf("expect token %q, actual: %q", testcase.token, added.Token())
				}
			})
		}

		t.Run("List", func(t *testing.T) {
			list, err := servers.List()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(list) != 2 {
				t.Fatalf("length mismatch: -want +got\n -%d\n +%d", 2, len(list))
			}
			first := list[0]
			if first.Host() != host1 {
				t.Errorf("expect host %q, actual: %q", host1, first.Host())
			}
			if first.User() != user1 {
				t.Errorf("expect user %q, actual: %q", user1, first.User())
			}
			if first.Token() != token1 {
				t.Errorf("expect token %q, actual: %q", token1, first.Token())
			}

			second := list[1]
			if second.Host() != host2 {
				t.Errorf("expect host %q, actual: %q", host2, second.Host())
			}
			if second.User() != user2 {
				t.Errorf("expect user %q, actual: %q", user2, second.User())
			}
			if second.Token() != token2 {
				t.Errorf("expect token %q, actual: %q", token2, second.Token())
			}
		})

		t.Run("SetDefault", func(t *testing.T) {
			if err := servers.SetDefault("unknown.dev"); !errors.Is(
				err,
				testtarget.ErrServerNotFound,
			) {
				t.Fatalf("expect error %q, actual: %v", testtarget.ErrServerNotFound, err)
			}
			if err := servers.SetDefault(host2); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			def, err := servers.Default()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if def.Host() != host2 {
				t.Errorf("expect host %q, actual: %q", host2, def.Host())
			}
			if def.User() != user2 {
				t.Errorf("expect user %q, actual: %q", user2, def.User())
			}
			if def.Token() != token2 {
				t.Errorf("expect token %q, actual: %q", token2, def.Token())
			}
		})

		t.Run("Remove", func(t *testing.T) {
			if err := servers.Remove("unknown.dev"); !errors.Is(err, testtarget.ErrServerNotFound) {
				t.Fatalf("expect error %q, actual: %v", testtarget.ErrServerNotFound, err)
			}
			if err := servers.Remove(host2); err != testtarget.ErrUnremovableServer {
				t.Fatalf("expect error %q, actual: %v", testtarget.ErrUnremovableServer, err)
			}

			if err := servers.Remove(host1); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, err := servers.Find(host1); !errors.Is(err, testtarget.ErrServerNotFound) {
				t.Fatalf("expect error %q, actual: %v", testtarget.ErrServerNotFound, err)
			}

			if err := servers.Remove(host2); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, err := servers.Find(host2); !errors.Is(err, testtarget.ErrNoServer) {
				t.Fatalf("expect error %q, actual: %v", testtarget.ErrNoServer, err)
			}
		})
	})
}
