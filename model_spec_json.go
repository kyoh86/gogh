package gogh

import "encoding/json"

type jsonSpec struct {
	Host  string `json:"host"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func (s *Spec) UnmarshalJSON(b []byte) error {
	var m jsonSpec
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	s.host = m.Host
	s.owner = m.Owner
	s.name = m.Name
	return nil
}

func (s Spec) MarshalJSON() ([]byte, error) {
	m := jsonSpec{
		Host:  s.host,
		Owner: s.owner,
		Name:  s.name,
	}
	return json.Marshal(m)
}

var (
	_ json.Unmarshaler = (*Spec)(nil)
	_ json.Marshaler   = (*Spec)(nil)
)
