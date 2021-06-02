package main

func defaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func defaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func defaultStringSlice(value, defaultValue []string) []string {
	if value == nil {
		return defaultValue
	}
	return value
}
