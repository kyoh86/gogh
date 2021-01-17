package alias

import "sync"

type Def struct {
	mutex   sync.Mutex
	lookup  lookup
	reverse reverse
}

// MarshalYAML implements the interface `yaml.Marshaler`
func (d *Def) MarshalYAML() (interface{}, error) {
	return d.lookup, nil
}

// UnmarshalYAML implements the interface `yaml.Unmarshaler`
func (d *Def) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed lookup
	if err := unmarshal(&parsed); err != nil {
		return err
	}
	d.lookup = parsed
	for alias, fullpath := range d.lookup {
		d.reverse.Set(fullpath, alias)
	}
	return nil
}

func (d *Def) Set(alias, fullpath string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.lookup.Has(alias) {
		d.reverse.Del(d.lookup.Get(alias), alias)
	}
	d.lookup.Set(alias, fullpath)
	d.reverse.Set(fullpath, alias)
}

func (d *Def) Del(alias string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !d.lookup.Has(alias) {
		return
	}
	fullpath := d.lookup.Get(alias)
	d.lookup.Del(alias)
	d.reverse.Del(fullpath, alias)
}

func (d *Def) Lookup(alias string) (fullpath string) {
	return d.lookup.Get(alias)
}

func (d *Def) Reverse(fullpath string) []string {
	return d.reverse.Get(fullpath).List()
}
