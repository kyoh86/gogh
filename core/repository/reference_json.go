package repository

import "encoding/json"

type referenceMarshalable struct {
	Host  string `json:"host"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func (r *Reference) UnmarshalJSON(b []byte) error {
	var m referenceMarshalable
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	r.host = m.Host
	r.owner = m.Owner
	r.name = m.Name
	return nil
}

func (r Reference) MarshalJSON() ([]byte, error) {
	m := referenceMarshalable{
		Host:  r.host,
		Owner: r.owner,
		Name:  r.name,
	}
	return json.Marshal(m)
}
