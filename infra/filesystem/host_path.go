package filesystem

import "github.com/kyoh86/gogh/v4/core/workspace"

func cloneHostPathAliases(aliases workspace.HostPathAliases) workspace.HostPathAliases {
	if len(aliases) == 0 {
		return nil
	}
	cloned := make(workspace.HostPathAliases, len(aliases))
	for host, alias := range aliases {
		if host == "" || alias == "" || host == alias {
			continue
		}
		cloned[host] = alias
	}
	if len(cloned) == 0 {
		return nil
	}
	return cloned
}

func reverseHostPathAliases(aliases workspace.HostPathAliases) map[string]string {
	if len(aliases) == 0 {
		return nil
	}
	reversed := make(map[string]string, len(aliases))
	for host, alias := range aliases {
		reversed[alias] = host
	}
	return reversed
}
