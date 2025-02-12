package main

var defaultFlag struct {
	BundleRestore bundleRestoreFlagsStruct `yaml:"bundleRestore,omitempty"`
	BundleDump    bundleDumpFlagsStruct    `yaml:"bundleDump,omitempty"`
	List          listFlagsStruct          `yaml:"list,omitempty"`
	Cwd           cwdFlagsStruct           `yaml:"cwd,omitempty"`
	Create        createFlagsStruct        `yaml:"create,omitempty"`
	Repos         reposFlagsStruct         `yaml:"repos,omitempty"`
	Fork          forkFlagsStruct          `yaml:"fork,omitempty"`
}
