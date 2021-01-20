package alias

// lookup is map[alias: string]fullpath: string
type lookup map[string]string

func (l *lookup) Set(key, value string) {
	if *l == nil {
		*l = map[string]string{}
	}
	(*l)[key] = value
}

func (l *lookup) Del(key string) {
	if *l == nil {
		return
	}
	delete(*l, key)
}

func (l *lookup) Get(key string) string {
	if *l == nil {
		return ""
	}
	return (*l)[key]
}

func (l *lookup) Has(key string) bool {
	if *l == nil {
		return false
	}
	_, ok := (*l)[key]
	return ok
}

// reverse is map[fullpath: string]Set[alias: string]
type reverse map[string]set

func (m *reverse) Set(key, value string) {
	if *m == nil {
		*m = map[string]set{}
	}
	c := (*m)[key]
	c.Set(value)
	(*m)[key] = c
}

func (m reverse) Del(key, value string) {
	if m == nil {
		return
	}
	c := m[key]
	c.Del(value)
	if len(c) == 0 {
		delete(m, key)
	}
}

func (m reverse) Get(key string) set {
	if m == nil {
		return nil
	}
	return m[key]
}

func (m reverse) Has(key, value string) bool {
	if m == nil {
		return false
	}
	c, ok := m[key]
	if !ok {
		return false
	}
	return c.Has(value)
}
