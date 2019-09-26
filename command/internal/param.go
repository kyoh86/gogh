package internal

func appendIf(array []string, flag string, value bool) []string {
	if value {
		return append(array, flag)
	}
	return array
}

func appendIfFilled(array []string, flag, value string) []string {
	if value == "" {
		return array
	}
	return append(array, flag, value)
}