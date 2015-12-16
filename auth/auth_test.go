package auth

import "testing"

func TestCheckScopes(t *testing.T) {
	{
		if err := checkScopes([]string{"public_repo", "repo", "user"}); err != nil {
			t.Error(err)
		}
	}
	{
		err := checkScopes([]string{}).(*scopeError)
		if err == nil {
			t.Errorf("Expect the err is not nil")
		}
		if !err.required["public_repo"] {
			t.Errorf("Expect 'public_repo' is required, but not required")
		}
		if !err.required["repo"] {
			t.Errorf("Expect 'repo' is required, but not required")
		}
		if !err.required["user"] {
			t.Errorf("Expect 'user' is required, but not required")
		}
	}
	{
		err := checkScopes(nil).(*scopeError)
		if err == nil {
			t.Errorf("Expect the err is not nil")
		}
		if !err.required["public_repo"] {
			t.Errorf("Expect 'public_repo' is required, but not required")
		}
		if !err.required["repo"] {
			t.Errorf("Expect 'repo' is required, but not required")
		}
		if !err.required["user"] {
			t.Errorf("Expect 'user' is required, but not required")
		}
	}
	{
		err := checkScopes([]string{"user", "admin"}).(*scopeError)
		if err == nil {
			t.Errorf("Expect the err is not nil")
		}
		if !err.required["public_repo"] {
			t.Errorf("Expect 'public_repo' is required, but not required")
		}
		if !err.required["repo"] {
			t.Errorf("Expect 'repo' is required, but not required")
		}
		if err.required["user"] {
			t.Errorf("Expect 'user' is not more required, but required")
		}
		if err.required["admin"] {
			t.Errorf("Expect 'admin' is not required, but required")
		}
	}
}
