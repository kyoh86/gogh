package util

// sets.go : helpers for arrays, slices and maps

// StringBoolMapKeys returns key strings of map[string]bool
func StringBoolMapKeys(m map[string]bool) (ret []string) {
	for k := range m {
		ret = append(ret, k)
	}
	return
}
