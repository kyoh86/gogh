package alias

// Lookup is map[alias: string]fullpath: string
type Lookup map[string]string

func (l *Lookup) Set(key, value string) {
	if *l == nil {
		*l = map[string]string{}
	}
	(*l)[key] = value
}

func (l *Lookup) Del(key string) {
	if *l == nil {
		return
	}
	delete(*l, key)
}

func (l *Lookup) Get(key string) string {
	if *l == nil {
		return ""
	}
	return (*l)[key]
}

func (l *Lookup) Has(key string) bool {
	if *l == nil {
		return false
	}
	_, ok := (*l)[key]
	return ok
}

func (l *Lookup) Keys() []string {
	if *l == nil {
		return nil
	}
	var keys []string
	for key := range *l {
		keys = append(keys, key)
	}
	return keys
}

// Reverse is map[fullpath: string]Set[alias: string]
type Reverse map[string]Set

func (m *Reverse) Set(key, value string) {
	if *m == nil {
		*m = map[string]Set{}
	}
	c := (*m)[key]
	c.Set(value)
	(*m)[key] = c
}

func (m Reverse) Del(key, value string) {
	if m == nil {
		return
	}
	c := m[key]
	c.Del(value)
	if len(c) == 0 {
		delete(m, key)
	}
}

func (m Reverse) Get(key string) Set {
	if m == nil {
		return nil
	}
	return m[key]
}

func (m Reverse) Has(key, value string) bool {
	if m == nil {
		return false
	}
	c, ok := m[key]
	if !ok {
		return false
	}
	return c.Has(value)
}

func (m Reverse) Keys() []string {
	if m == nil {
		return nil
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
