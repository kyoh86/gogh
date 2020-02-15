package config

func MergeConfig(base *Config, override ...*Config) *Config {
	c := *base
	for _, o := range override {
		c.VRoot = mergePathListOption(c.VRoot, o.VRoot)
		c.GitHub.Token = mergeStringOption(c.GitHub.Token, o.GitHub.Token)
		c.GitHub.User = mergeStringOption(c.GitHub.User, o.GitHub.User)
		c.GitHub.Host = mergeStringOption(c.GitHub.Host, o.GitHub.Host)
	}
	return &c
}

func mergeBoolOption(base, override BoolOption) BoolOption {
	switch {
	case override != EmptyBoolOption:
		return override
	case base != EmptyBoolOption:
		return base
	default:
		return EmptyBoolOption
	}
}

func mergeStringOption(base, override string) string {
	if override != "" {
		return override
	}
	return base
}

func mergePathListOption(base, override []string) []string {
	if len(override) > 0 {
		return override
	}
	return base
}
