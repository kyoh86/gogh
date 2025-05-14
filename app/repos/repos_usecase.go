package repos

import (
	"context"
	"fmt"
	"iter"

	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/typ"
)

type UseCase struct {
	hostingService hosting.HostingService
}

func NewUseCase(hostingService hosting.HostingService) *UseCase {
	return &UseCase{
		hostingService: hostingService,
	}
}

type Options struct {
	Limit    int
	Privacy  string
	Fork     string
	Archive  string
	Format   string
	Color    string
	Relation []string
	Sort     string
	Order    string
}

func convertOpts(f Options) (hosting.ListRepositoryOptions, error) {
	var opts hosting.ListRepositoryOptions

	// Set limit
	switch f.Limit {
	case 0:
		opts.Limit = 30
	case -1:
		opts.Limit = 0 // no limit
	default:
		opts.Limit = f.Limit
	}

	if err := typ.Remap(&opts.Privacy, map[string]hosting.RepositoryPrivacy{
		"private": hosting.RepositoryPrivacyPrivate,
		"public":  hosting.RepositoryPrivacyPublic,
	}, f.Privacy); err != nil {
		return opts, fmt.Errorf("invalid privacy option %q", f.Privacy)
	}
	if err := typ.Remap(&opts.IsFork, map[string]hosting.Tristate{
		"forked":     hosting.TristateTrue,
		"not-forked": hosting.TristateFalse,
	}, f.Fork); err != nil {
		return opts, fmt.Errorf("invalid fork option %q", f.Fork)
	}
	if err := typ.Remap(&opts.IsArchived, map[string]hosting.Tristate{
		"archived":     hosting.TristateTrue,
		"not-archived": hosting.TristateFalse,
	}, f.Archive); err != nil {
		return opts, fmt.Errorf("invalid archive option %q", f.Archive)
	}

	// Set relation filters
	for _, r := range f.Relation {
		if field, exists := map[string]hosting.RepositoryAffiliation{
			"owner":               hosting.RepositoryAffiliationOwner,
			"organization-member": hosting.RepositoryAffiliationOrganizationMember,
			"collaborator":        hosting.RepositoryAffiliationCollaborator,
		}[r]; exists {
			opts.OwnerAffiliations = append(opts.OwnerAffiliations, field)
		} else {
			return opts, fmt.Errorf("invalid relation %q", r)
		}
	}

	if err := typ.Remap(&opts.OrderBy.Field, map[string]hosting.RepositoryOrderField{
		"created-at": hosting.RepositoryOrderFieldCreatedAt,
		"name":       hosting.RepositoryOrderFieldName,
		"pushed-at":  hosting.RepositoryOrderFieldPushedAt,
		"stargazers": hosting.RepositoryOrderFieldStargazers,
		"updated-at": hosting.RepositoryOrderFieldUpdatedAt,
	}, f.Sort); err != nil {
		return opts, fmt.Errorf("invalid sort field %q", f.Sort)
	}

	if err := typ.Remap(&opts.OrderBy.Direction, map[string]hosting.OrderDirection{
		"asc":        hosting.OrderDirectionAsc,
		"ascending":  hosting.OrderDirectionAsc,
		"desc":       hosting.OrderDirectionDesc,
		"descending": hosting.OrderDirectionDesc,
	}, f.Order); err != nil {
		return opts, fmt.Errorf("invalid order direction %q", f.Order)
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
