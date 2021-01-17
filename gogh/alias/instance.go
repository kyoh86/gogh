package alias

var Instance Def

func Set(alias, fullpath string) {
	Instance.Set(alias, fullpath)
}

func Del(alias string) {
	Instance.Del(alias)
}

func Lookup(alias string) (fullpath string) {
	return Instance.Lookup(alias)
}

func Reverse(fullpath string) []string {
	return Instance.Reverse(fullpath)
}
