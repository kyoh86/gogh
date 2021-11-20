package main

var config struct {
	Roots []expandedPath `yaml:"roots"`
}

func defaultRoot() string {
	return config.Roots[0].expanded
}

func roots() []string {
	list := make([]string, 0, len(config.Roots))
	for _, r := range config.Roots {
		list = append(list, r.expanded)
	}
	return list
}

func setDefaultRoot(r string) error {
	rootList := make([]expandedPath, 0, len(config.Roots))
	newDefault, err := parsePath(r)
	if err != nil {
		return err
	}
	rootList = append(rootList, newDefault)
	for _, root := range config.Roots {
		if root.raw == r {
			continue
		}
		rootList = append(rootList, root)
	}
	config.Roots = rootList
	return nil
}

func addRoots(rootList []string) error {
	for _, r := range rootList {
		newRoot, err := parsePath(r)
		if err != nil {
			return err
		}
		config.Roots = append(config.Roots, newRoot)
	}
	return nil
}

func removeRoot(r string) {
	rootList := make([]expandedPath, 0, len(config.Roots))
	for _, root := range config.Roots {
		if root.raw == r || root.expanded == r {
			continue
		}
		rootList = append(rootList, root)
	}
	config.Roots = rootList
}
