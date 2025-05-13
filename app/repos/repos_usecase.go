package repos

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"strings"

	"github.com/kyoh86/gogh/v3/core/hosting"
)

type UseCase struct {
	hostingService hosting.HostingService
}

func NewUseCase(hostingService hosting.HostingService) *UseCase {
	return &UseCase{
		hostingService: hostingService,
	}
}

// TODO: app/config/flags.go? こいつだけ場所が変
type Options struct {
	Limit       int      `yaml:"limit,omitempty" toml:"limit,omitempty"`
	Public      bool     `yaml:"public,omitempty" toml:"public,omitempty"`
	Private     bool     `yaml:"private,omitempty" toml:"private,omitempty"`
	Fork        bool     `yaml:"fork,omitempty" toml:"fork,omitempty"`
	NotFork     bool     `yaml:"notFork,omitempty" toml:"notFork,omitempty"`
	Archived    bool     `yaml:"archived,omitempty" toml:"archived,omitempty"`
	NotArchived bool     `yaml:"notArchived,omitempty" toml:"notArchived,omitempty"`
	Format      string   `yaml:"format,omitempty" toml:"format,omitempty"`
	Color       string   `yaml:"color,omitempty" toml:"color,omitempty"`
	Relation    []string `yaml:"relation,omitempty" toml:"relation,omitempty"`
	Sort        string   `yaml:"sort,omitempty" toml:"sort,omitempty"`
	Order       string   `yaml:"order,omitempty" toml:"order,omitempty"`
}

var relationMap = map[string]hosting.RepositoryAffiliation{
	"owner":               hosting.RepositoryAffiliationOwner,
	"organization-member": hosting.RepositoryAffiliationOrganizationMember,
	"collaborator":        hosting.RepositoryAffiliationCollaborator,
}
var RelationAccepts = []string{
	"owner",
	"organization-member",
	"collaborator",
}
var sortMap = map[string]hosting.RepositoryOrderField{
	"created-at": hosting.RepositoryOrderFieldCreatedAt,
	"name":       hosting.RepositoryOrderFieldName,
	"pushed-at":  hosting.RepositoryOrderFieldPushedAt,
	"stargazers": hosting.RepositoryOrderFieldStargazers,
	"updated-at": hosting.RepositoryOrderFieldUpdatedAt,
}
var orderMap = map[string]hosting.OrderDirection{
	"asc":        hosting.OrderDirectionAsc,
	"ascending":  hosting.OrderDirectionAsc,
	"desc":       hosting.OrderDirectionDesc,
	"descending": hosting.OrderDirectionDesc,
}

func convertOpts(f Options) (hosting.ListRepositoryOptions, error) {
	var opts hosting.ListRepositoryOptions
	// Check mutually exclusive flags
	if f.Private && f.Public {
		return opts, errors.New("specify only one of `--private` or `--public`")
	}
	if f.Fork && f.NotFork {
		return opts, errors.New("specify only one of `--fork` or `--no-fork`")
	}
	if f.Archived && f.NotArchived {
		return opts, errors.New("specify only one of `--archived` or `--no-archived`")
	}

	// Set limit
	switch f.Limit {
	case 0:
		opts.Limit = 30
	case -1:
		opts.Limit = 0 // no limit
	default:
		opts.Limit = f.Limit
	}

	// Set privacy filter
	if f.Private {
		opts.Privacy = hosting.RepositoryPrivacyPrivate
	}
	if f.Public {
		opts.Privacy = hosting.RepositoryPrivacyPublic
	}

	// Set fork and archive filters
	if f.Fork {
		opts.IsFork = hosting.BooleanFilterTrue
	}
	if f.NotFork {
		opts.IsFork = hosting.BooleanFilterFalse
	}
	if f.Archived {
		opts.IsArchived = hosting.BooleanFilterTrue
	}
	if f.NotArchived {
		opts.IsArchived = hosting.BooleanFilterFalse
	}

	// Set relation filters
	for _, r := range f.Relation {
		if field, exists := relationMap[r]; exists {
			opts.OwnerAffiliations = append(opts.OwnerAffiliations, field)
		} else {
			return opts, fmt.Errorf("invalid relation %q", r)
		}
	}

	// Set sort field
	if field, exists := sortMap[strings.ToLower(f.Sort)]; exists {
		opts.OrderBy.Field = field
	} else {
		return opts, fmt.Errorf("invalid sort field %q", f.Sort)
	}

	// Set sort direction
	if field, exists := orderMap[strings.ToLower(f.Order)]; exists {
		opts.OrderBy.Direction = field
	} else {
		return opts, fmt.Errorf("invalid order field %q", f.Order)
	}

	return opts, nil
}

func (uc *UseCase) Execute(ctx context.Context, opts Options) iter.Seq2[*hosting.Repository, error] {
	return func(yield func(*hosting.Repository, error) bool) {
		listOpts, err := convertOpts(opts)
		if err != nil {
			yield(nil, err)
			return
		}
		for repo, err := range uc.hostingService.ListRepository(ctx, listOpts) {
			if !yield(repo, err) {
				return
			}
		}
	}
}
