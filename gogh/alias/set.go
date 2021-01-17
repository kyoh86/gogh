package alias

type Set map[string]struct{}

func NewSet(list ...string) Set {
	var s Set
	for _, item := range list {
		s.Set(item)
	}
	return s
}

func (s *Set) Set(key string) {
	if *s == nil {
		*s = map[string]struct{}{}
	}
	(*s)[key] = struct{}{}
}

func (s Set) Del(key string) {
	if s == nil {
		return
	}
	delete(s, key)
}

func (s Set) Has(key string) bool {
	if s == nil {
		return false
	}
	_, ok := s[key]
	return ok
}

func (s Set) List() []string {
	var keys []string
	for key := range s {
		keys = append(keys, key)
	}
	return keys
}
