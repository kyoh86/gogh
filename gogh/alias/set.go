package alias

type set map[string]struct{}

func newSet(list ...string) set {
	var s set
	for _, item := range list {
		s.Set(item)
	}
	return s
}

func (s *set) Set(key string) {
	if *s == nil {
		*s = map[string]struct{}{}
	}
	(*s)[key] = struct{}{}
}

func (s set) Del(key string) {
	if s == nil {
		return
	}
	delete(s, key)
}

func (s set) Has(key string) bool {
	if s == nil {
		return false
	}
	_, ok := s[key]
	return ok
}

func (s set) List() []string {
	keys := make([]string, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	return keys
}
