package gogh

// falsePtr converts bool (default: false) to *bool (default: nil as true)
func falsePtr(b bool) *bool {
	if b {
		f := false
		return &f
	}
	return nil // == &true
}

func boolPtr(b bool) *bool {
	if b == false {
		return nil
	}
	return &b
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}
