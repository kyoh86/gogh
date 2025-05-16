package typ

type Map[TKey comparable, TVal any] map[TKey]TVal

func (m *Map[TKey, TVal]) Set(key TKey, val TVal) {
	if *m == nil {
		*m = map[TKey]TVal{}
	}
	(*m)[key] = val
}

func (m *Map[TKey, TVal]) Delete(key TKey) {
	if *m == nil {
		return
	}
	delete(*m, key)
}

func (m *Map[TKey, TVal]) Has(key TKey) bool {
	if *m == nil {
		return false
	}
	_, ok := (*m)[key]
	return ok
}

func (m *Map[TKey, TVal]) TryGet(key TKey) (TVal, bool) {
	var d TVal
	if *m == nil {
		return d, false
	}
	v, ok := (*m)[key]
	return v, ok
}

func (m *Map[TKey, TVal]) GetOrSet(key TKey, setValue TVal) TVal {
	if *m == nil {
		*m = map[TKey]TVal{
			key: setValue,
		}
		return setValue
	}
	if v, ok := (*m)[key]; ok {
		return v
	}
	(*m)[key] = setValue
	return setValue
}
