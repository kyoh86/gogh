package gogh

func unique(items []string) (uniq []string) {
	dups := map[string]struct{}{}

	for _, item := range items {
		if _, ok := dups[item]; ok {
			continue
		}
		dups[item] = struct{}{}
		uniq = append(uniq, item)
	}
	return
}
