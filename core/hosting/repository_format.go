package hosting

import (
	"encoding/json"
	"time"
)

// RepositoryFormat is an interface that defines a method to format a repository
type RepositoryFormat interface {
	Format(r Repository) (string, error)
}

// RepositoryFormatFunc is a function type that implements the RepositoryFormat interface
type RepositoryFormatFunc func(Repository) (string, error)

// Format calls the function itself to format the repository
func (f RepositoryFormatFunc) Format(r Repository) (string, error) {
	return f(r)
}

// RepositoryFormatRef formats the repository reference as a string
var RepositoryFormatRef = RepositoryFormatFunc(func(r Repository) (string, error) {
	return r.Ref.String(), nil
})

// RepositoryFormatURL formats the repository URL as a string
var RepositoryFormatURL = RepositoryFormatFunc(func(r Repository) (string, error) {
	return r.URL, nil
})

type jsonRef struct {
	Host  string `json:"host"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type jsonParent struct {
	Ref      jsonRef `json:"ref"`
	CloneURL string  `json:"cloneUrl,omitempty"`
}

// RepositoryFormatJSON formats the repository as a JSON string in oneline
var RepositoryFormatJSON = RepositoryFormatFunc(func(r Repository) (string, error) {
	j := struct {
		Ref         jsonRef     `json:"ref"`
		URL         string      `json:"url"`
		CloneURL    string      `json:"cloneUrl,omitempty"`
		UpdatedAt   time.Time   `json:"updatedAt"`
		Parent      *jsonParent `json:"parent,omitempty"`
		Description string      `json:"description,omitempty"`
		Homepage    string      `json:"homepage,omitempty"`
		Language    string      `json:"language,omitempty"`
		Archived    bool        `json:"archived,omitempty"`
		Private     bool        `json:"private,omitempty"`
		IsTemplate  bool        `json:"isTemplate,omitempty"`
		Fork        bool        `json:"fork,omitempty"`
	}{
		Ref: jsonRef{
			Host:  r.Ref.Host(),
			Owner: r.Ref.Owner(),
			Name:  r.Ref.Name(),
		},
		URL:         r.URL,
		CloneURL:    r.CloneURL,
		UpdatedAt:   r.UpdatedAt,
		Description: r.Description,
		Homepage:    r.Homepage,
		Language:    r.Language,
		Archived:    r.Archived,
		Private:     r.Private,
		IsTemplate:  r.IsTemplate,
		Fork:        r.Fork,
	}
	if r.Parent != nil {
		j.Parent = &jsonParent{
			Ref: jsonRef{
				Host:  r.Parent.Ref.Host(),
				Owner: r.Parent.Ref.Owner(),
				Name:  r.Parent.Ref.Name(),
			},
			CloneURL: r.Parent.CloneURL,
		}
	}
	buf, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(buf), nil
})
