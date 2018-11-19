package gogh

import "context"

type mockContext struct {
	context.Context
	userName      string
	userNameError error
	roots         []string
	rootsError    error
	gheHosts      []string
	gheHostsError error
}

func (c *mockContext) UserName() (string, error) {
	return c.userName, c.userNameError
}

func (c *mockContext) Roots() ([]string, error) {
	return c.roots, c.rootsError
}

func (c *mockContext) PrimaryRoot() (string, error) {
	rts, err := c.Roots()
	if err != nil {
		return "", err
	}
	return rts[0], nil
}

func (c *mockContext) GHEHosts() ([]string, error) {
	return c.gheHosts, c.gheHostsError
}
