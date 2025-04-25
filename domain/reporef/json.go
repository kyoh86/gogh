package reporef

import "encoding/json"

type jsonRepoRef struct {
	Host  string `json:"host"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func (r *RepoRef) UnmarshalJSON(b []byte) error {
	var m jsonRepoRef
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	r.host = m.Host
	r.owner = m.Owner
	r.name = m.Name
	return nil
}

func (r RepoRef) MarshalJSON() ([]byte, error) {
	m := jsonRepoRef{
		Host:  r.host,
		Owner: r.owner,
		Name:  r.name,
	}
	return json.Marshal(m)
}

var (
	_ json.Unmarshaler = (*RepoRef)(nil)
	_ json.Marshaler   = (*RepoRef)(nil)
)
