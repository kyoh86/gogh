package delegate

func AppendIf(array []string, flag string, value bool) []string {
	if value {
		return append(array, flag)
	}
	return array
}

func AppendIfFilled(array []string, flag, value string) []string {
	if value == "" {
		return array
	}
	return append(array, flag, value)
}

func AppendPairedIfFilled(array []string, flag, value string) []string {
	if value == "" {
		return array
	}
	return append(array, flag+"="+value)
}
