package github

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"net/url"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/google/go-github/v69/github"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/infra/githubv4"
	"github.com/kyoh86/gogh/v4/typ"
	"golang.org/x/oauth2"
)

type HostingService struct {
	// GitHub client, etc.
	tokenService       auth.TokenService
	defaultNameService repository.DefaultNameService
	knownOwners        map[string]string
}

const (
	GlobalHost    = "github.com"
	GlobalAPIHost = "api.github.com"

	ClientID = "Ov23li6aEWIxek6F8P5L"
)

func oAuth2Config(host string) *oauth2.Config {
	return &oauth2.Config{
		ClientID: ClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:       fmt.Sprintf("https://%s/login/oauth/authorize", host),
			TokenURL:      fmt.Sprintf("https://%s/login/oauth/access_token", host),
			DeviceAuthURL: fmt.Sprintf("https://%s/login/device/code", host),
		},
		Scopes: []string{string(github.ScopeRepo), string(github.ScopeDeleteRepo)},
	}
}

type tokenSource struct {
	ctx   context.Context
	host  string
	token *oauth2.Token
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	if s.token.Valid() {
		return s.token, nil
	}
	newToken, err := refreshAccessToken(s.ctx, s.host, s.token)
	if err != nil {
		return nil, err
	}
	s.token = newToken
	return newToken, nil
}

func refreshAccessToken(ctx context.Context, host string, token *oauth2.Token) (*oauth2.Token, error) {
	oauthConfig := oAuth2Config(host)
	tokenSource := oauthConfig.TokenSource(ctx, &oauth2.Token{RefreshToken: token.RefreshToken})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func getSource(ctx context.Context, host string, token *auth.Token) oauth2.TokenSource {
	if token == nil {
		return nil
	}
	return oauth2.ReuseTokenSource(token, &tokenSource{ctx: ctx, host: host, token: token})
}

type connection struct {
	rest *github.Client
	gql  graphql.Client
}

func getEnterpriseConnection(host string, httpClient *http.Client) *connection {
	baseRESTURL := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/v3",
	}
	uploadRESTURL := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/uploads",
	}
	baseGQLURL := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/graphql",
	}
	restClient, err := github.NewClient(httpClient).WithEnterpriseURLs(baseRESTURL.String(), uploadRESTURL.String())
	if err != nil {
		// NOTE: WithEnterpriseURLs returns error if the URLs are not valid.
		// We assume that the URLs are valid.
		panic(err)
	}
	return &connection{rest: restClient, gql: graphql.NewClient(baseGQLURL.String(), httpClient)}
}

func getConnection(ctx context.Context, host string, token *auth.Token) *connection {
	source := getSource(ctx, host, token)
	httpClient := oauth2.NewClient(ctx, source)
	if host == GlobalHost || host == GlobalAPIHost {
		return &connection{
			rest: github.NewClient(httpClient),
			gql:  graphql.NewClient("https://"+GlobalAPIHost+"/graphql", httpClient),
		}
	}
	return getEnterpriseConnection(host, httpClient)
}

// NewHostingService creates a new HostingService instance
func NewHostingService(tokenService auth.TokenService, defaultNameService repository.DefaultNameService) *HostingService {
	return &HostingService{
		tokenService:       tokenService,
		defaultNameService: defaultNameService,
		knownOwners:        map[string]string{},
	}
}

// GetURLOf implements hosting.HostingService.
func (s *HostingService) GetURLOf(ref repository.Reference) (*url.URL, error) {
	return &url.URL{
		Scheme: "https",
		Host:   ref.Host(),
		Path:   strings.Join([]string{ref.Owner(), ref.Name()}, "/"),
	}, nil
}

// ParseURL implements hosting.HostingService.
func (s *HostingService) ParseURL(u *url.URL) (*repository.Reference, error) {
	words := strings.SplitN(strings.TrimPrefix(strings.TrimSuffix(u.Path, ".git"), "/"), "/", 2)
	if len(words) < 2 {
		return nil, fmt.Errorf("invalid path: %q", u.Path)
	}
	return typ.Ptr(repository.NewReference(u.Host, words[0], strings.TrimSuffix(words[1], ".git"))), nil
}

var ErrTokenNotFound = errors.New("no token found")

// GetTokenFor cache requested token for the host and owner
func (s *HostingService) GetTokenFor(ctx context.Context, host, owner string) (string, auth.Token, error) {
	key := strings.Join([]string{host, owner}, "/")
	tokenOwner, ok := s.knownOwners[key]
	if ok {
		_, token, err := s.getTokenForCore(ctx, host, tokenOwner)
		return tokenOwner, token, err
	}
	tokenOwner, token, err := s.getTokenForCore(ctx, host, owner)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			// If no token is found, use the default owner as the username
			defaultOwner, inErr := s.defaultNameService.GetDefaultOwnerFor(host)
			if inErr != nil {
				return "", token, fmt.Errorf("getting default owner: %w", inErr)
			}
			tokenOwner, token, err = s.getTokenForCore(ctx, host, defaultOwner)
			if err != nil {
				return "", token, fmt.Errorf("getting default token: %w", inErr)
			} else {
				s.knownOwners[key] = tokenOwner
			}
		} else {
			return "", token, fmt.Errorf("getting token for %s/%s: %w", host, owner, err)
		}
	} else {
		s.knownOwners[key] = tokenOwner
	}
	return tokenOwner, token, err
}

func memberOf(ctx context.Context, owner string, entry auth.TokenEntry) bool {
	connection := getConnection(ctx, entry.Host, &entry.Token)
	orgs, _, err := connection.rest.Organizations.List(ctx, "", &github.ListOptions{PerPage: 100})
	if err != nil {
		return false
	}
	for _, o := range orgs {
		if *o.Login == owner {
			return true
		}
	}
	return false
}

func (s *HostingService) getTokenForCore(ctx context.Context, host, owner string) (string, auth.Token, error) {
	if s.tokenService.Has(host, owner) {
		token, err := s.tokenService.Get(host, owner)
		return owner, token, err
	}

	for _, entry := range s.tokenService.Entries() {
		if entry.Host != host {
			continue
		}
		if entry.Owner == owner {
			return entry.Owner, entry.Token, nil
		}

		// Check if this user is a member of the target organization
		if memberOf(ctx, owner, entry) {
			return entry.Owner, entry.Token, nil
		}
	}

	return "", auth.Token{}, ErrTokenNotFound
}

// GetRepository retrieves repository information from a remote source
func (s *HostingService) GetRepository(ctx context.Context, reference repository.Reference) (*hosting.Repository, error) {
	_, token, err := s.GetTokenFor(ctx, reference.Host(), reference.Owner())
	if err != nil {
		return nil, fmt.Errorf("getting token for %s/%s: %w", reference.Host(), reference.Owner(), err)
	}
	conn := getConnection(ctx, reference.Host(), &token)
	ghRepo, _, err := conn.rest.Repositories.Get(ctx, reference.Owner(), reference.Name())
	if err != nil {
		return nil, fmt.Errorf("requesting repository: %w", err)
	}
	// Convert github.Repository to hosting.Repository
	repo := &hosting.Repository{
		Ref:         reference,
		Description: ghRepo.GetDescription(),
		Homepage:    ghRepo.GetHomepage(),
		Language:    ghRepo.GetLanguage(),
		Fork:        ghRepo.GetFork(),
		Archived:    ghRepo.GetArchived(),
		Private:     ghRepo.GetPrivate(),
		URL:         ghRepo.GetHTMLURL(),
		CloneURL:    ghRepo.GetCloneURL(),
		UpdatedAt:   ghRepo.GetUpdatedAt().Time,
	}

	// Set if the parent repository is available
	if parent := ghRepo.GetParent(); parent != nil {
		u, err := url.Parse(parent.GetHTMLURL())
		if err != nil {
			return nil, fmt.Errorf("invalid parent HTML URL: %w", err)
		}
		parentRepo := &hosting.ParentRepository{
			Ref:      repository.NewReference(u.Host, parent.GetOwner().GetLogin(), parent.GetName()),
			CloneURL: ghRepo.GetParent().GetCloneURL(),
		}
		repo.Parent = parentRepo
	}
	return repo, err
}

const (
	RepoListMaxLimitPerPage = 100
)

// ListRepository retrieves a list of repositories from a remote source
func (s *HostingService) ListRepository(ctx context.Context, opts hosting.ListRepositoryOptions) iter.Seq2[*hosting.Repository, error] {
	return func(yield func(*hosting.Repository, error) bool) {
		var limit int
		switch {
		case opts.Limit == 0:
			limit = RepoListMaxLimitPerPage
		case opts.Limit > RepoListMaxLimitPerPage:
			limit = RepoListMaxLimitPerPage
		default:
			limit = opts.Limit
		}
		var count int
		for _, entry := range s.tokenService.Entries() {
			if err := ctx.Err(); err != nil {
				yield(nil, err)
				return
			}

			conn := getConnection(ctx, entry.Host, &entry.Token)
			var after string

			for {
				if err := ctx.Err(); err != nil {
					yield(nil, err)
					return
				}

				var privacy githubv4.RepositoryPrivacy
				if err := typ.Remap(&privacy, map[hosting.RepositoryPrivacy]githubv4.RepositoryPrivacy{
					hosting.RepositoryPrivacyPublic:  githubv4.RepositoryPrivacyPublic,
					hosting.RepositoryPrivacyPrivate: githubv4.RepositoryPrivacyPrivate,
				}, opts.Privacy); err != nil {
					yield(nil, fmt.Errorf("invalid privacy option %q", opts.Privacy))
					return
				}
				orderField := githubv4.RepositoryOrderFieldUpdatedAt
				if err := typ.Remap(&orderField, map[hosting.RepositoryOrderField]githubv4.RepositoryOrderField{
					hosting.RepositoryOrderFieldCreatedAt:  githubv4.RepositoryOrderFieldCreatedAt,
					hosting.RepositoryOrderFieldUpdatedAt:  githubv4.RepositoryOrderFieldUpdatedAt,
					hosting.RepositoryOrderFieldPushedAt:   githubv4.RepositoryOrderFieldPushedAt,
					hosting.RepositoryOrderFieldName:       githubv4.RepositoryOrderFieldName,
					hosting.RepositoryOrderFieldStargazers: githubv4.RepositoryOrderFieldStargazers,
				}, opts.OrderBy.Field); err != nil {
					yield(nil, fmt.Errorf("invalid order field %q", opts.OrderBy.Field))
					return
				}
				orderDirection := githubv4.OrderDirectionDesc
				if err := typ.Remap(&orderDirection, map[hosting.OrderDirection]githubv4.OrderDirection{
					hosting.OrderDirectionAsc:  githubv4.OrderDirectionAsc,
					hosting.OrderDirectionDesc: githubv4.OrderDirectionDesc,
				}, opts.OrderBy.Direction); err != nil {
					yield(nil, fmt.Errorf("invalid order direction %q", opts.OrderBy.Direction))
					return
				}
				affs, err := convertOwnerAffiliations(opts.OwnerAffiliations)
				if err != nil {
					yield(nil, fmt.Errorf("invalid owner affiliations %q: %w", opts.OwnerAffiliations, err))
					return
				}

				isFork, err := opts.IsFork.AsBoolPtr()
				if err != nil {
					yield(nil, fmt.Errorf("invalid isFork option %q: %w", opts.IsFork, err))
					return
				}
				isArchived, err := opts.IsArchived.AsBoolPtr()
				if err != nil {
					yield(nil, fmt.Errorf("invalid isArchived option %q: %w", opts.IsArchived, err))
					return
				}
				repos, err := githubv4.ListRepos(
					ctx,
					conn.gql,
					limit,
					after,
					isFork,
					privacy,
					affs,
					githubv4.RepositoryOrder{Field: orderField, Direction: orderDirection},
					isArchived,
				)
				if err != nil {
					yield(nil, err)
					return
				}

				repositories := repos.Viewer.Repositories
				for _, edge := range repositories.Edges {
					if err := ctx.Err(); err != nil {
						yield(nil, err)
						return
					}
					if !yield(typ.Ptr(convertRepositoryFragment(entry.Host, edge.Node.RepositoryFragment)), nil) {
						return
					}

					count++
					if opts.Limit > 0 && count >= opts.Limit {
						return
					}
				}

				page := repositories.PageInfo.PageInfoFragment
				if !page.HasNextPage {
					break
				}
				after = page.EndCursor
			}

			if opts.Limit > 0 && count >= opts.Limit {
				return
			}
		}
	}
}

// invertPtr returns nil if b is true, or pointer to false if b is false.
func invertPtr(b bool) *bool {
	if b {
		f := false
		return &f
	}
	return nil // true returns nil
}

func (s *HostingService) CreateRepository(
	ctx context.Context,
	ref repository.Reference,
	opts hosting.CreateRepositoryOptions,
) (*hosting.Repository, error) {
	user, token, err := s.GetTokenFor(ctx, ref.Host(), ref.Owner())
	if err != nil {
		return nil, fmt.Errorf("getting token for %s/%s: %w", ref.Host(), ref.Owner(), err)
	}
	conn := getConnection(ctx, ref.Host(), &token)
	org := ""
	if user != ref.Owner() {
		org = ref.Owner()
	}
	repo, _, err := conn.rest.Repositories.Create(ctx, org, &github.Repository{
		Name:                typ.NilablePtr(ref.Name()),
		Description:         typ.NilablePtr(opts.Description),
		Homepage:            typ.NilablePtr(opts.Homepage),
		Private:             typ.NilablePtr(opts.Private),
		HasIssues:           invertPtr(opts.DisableIssues),
		HasProjects:         invertPtr(opts.DisableProjects),
		HasWiki:             invertPtr(opts.DisableWiki),
		HasDownloads:        invertPtr(opts.DisableDownloads),
		IsTemplate:          typ.NilablePtr(opts.IsTemplate),
		TeamID:              typ.NilablePtr(opts.TeamID),
		AutoInit:            typ.NilablePtr(opts.AutoInit),
		GitignoreTemplate:   typ.NilablePtr(opts.GitignoreTemplate),
		LicenseTemplate:     typ.NilablePtr(opts.LicenseTemplate),
		AllowSquashMerge:    invertPtr(opts.PreventSquashMerge),
		AllowMergeCommit:    invertPtr(opts.PreventMergeCommit),
		AllowRebaseMerge:    invertPtr(opts.PreventRebaseMerge),
		DeleteBranchOnMerge: typ.NilablePtr(opts.DeleteBranchOnMerge),
	})
	if err != nil {
		return nil, fmt.Errorf("requesting new repository: %w", err)
	}
	return convertRepository(ref, repo)
}

func (s *HostingService) CreateRepositoryFromTemplate(
	ctx context.Context,
	ref repository.Reference,
	tmp repository.Reference,
	opts hosting.CreateRepositoryFromTemplateOptions,
) (*hosting.Repository, error) {
	user, token, err := s.GetTokenFor(ctx, ref.Host(), ref.Owner())
	if err != nil {
		return nil, fmt.Errorf("getting token for %s/%s: %w", ref.Host(), ref.Owner(), err)
	}
	conn := getConnection(ctx, ref.Host(), &token)
	req := github.TemplateRepoRequest{
		Name:               typ.Ptr(ref.Name()),
		Description:        &opts.Description,
		IncludeAllBranches: &opts.IncludeAllBranches,
		Private:            &opts.Private,
	}
	if user != ref.Owner() {
		req.Owner = typ.Ptr(ref.Owner())
	}
	repo, _, err := conn.rest.Repositories.CreateFromTemplate(
		ctx,
		tmp.Owner(),
		tmp.Name(),
		&req,
	)
	if err != nil {
		return nil, fmt.Errorf("requesting new repository from template: %w", err)
	}
	return convertRepository(ref, repo)
}

// DeleteRepository deletes a repository from a remote source
func (s *HostingService) DeleteRepository(ctx context.Context, reference repository.Reference) error {
	_, token, err := s.GetTokenFor(ctx, reference.Host(), reference.Owner())
	if err != nil {
		return fmt.Errorf("getting token for %s/%s: %w", reference.Host(), reference.Owner(), err)
	}
	conn := getConnection(ctx, reference.Host(), &token)
	if _, err = conn.rest.Repositories.Delete(ctx, reference.Owner(), reference.Name()); err != nil {
		return fmt.Errorf("requesting to delete repository: %w", err)
	}
	return nil
}

// ForkRepository implements hosting.HostingService.
func (s *HostingService) ForkRepository(
	ctx context.Context,
	ref repository.Reference,
	target repository.Reference,
	opts hosting.ForkRepositoryOptions,
) (*hosting.Repository, error) {
	user, token, err := s.GetTokenFor(ctx, target.Host(), target.Owner())
	if err != nil {
		return nil, fmt.Errorf("getting token for %s/%s: %w", ref.Host(), ref.Owner(), err)
	}
	conn := getConnection(ctx, ref.Host(), &token)
	ghOpts := &github.RepositoryCreateForkOptions{
		Name:              target.Name(),
		DefaultBranchOnly: opts.DefaultBranchOnly,
	}
	if user != target.Owner() {
		ghOpts.Organization = target.Owner()
	}
	fork, _, err := conn.rest.Repositories.CreateFork(ctx, ref.Owner(), ref.Name(), ghOpts)
	if err != nil {
		var acc *github.AcceptedError
		if errors.As(err, &acc) {
			got, _, err := conn.rest.Repositories.Get(ctx, target.Owner(), ref.Name())
			if err != nil {
				return nil, fmt.Errorf("requesting forked repository: %w", err)
			}
			fork = got
		} else {
			return nil, fmt.Errorf("requesting fork: %w", err)
		}
	}
	return convertRepository(target, fork)
}

func convertRepository(target repository.Reference, repo *github.Repository) (*hosting.Repository, error) {
	var parent *hosting.ParentRepository
	if raw := repo.GetParent(); raw != nil {
		u, err := url.Parse(raw.GetHTMLURL())
		if err != nil {
			return nil, fmt.Errorf("invalid HTML URL: %w", err)
		}
		ref := repository.NewReference(
			u.Host,
			raw.GetOwner().GetLogin(),
			raw.GetName(),
		)
		parent = &hosting.ParentRepository{
			CloneURL: raw.GetHTMLURL(),
			Ref:      ref,
		}
	}
	return &hosting.Repository{
		Ref:         target,
		URL:         repo.GetHTMLURL(),
		Parent:      parent,
		CloneURL:    repo.GetCloneURL(),
		Description: repo.GetDescription(),
		Homepage:    repo.GetHomepage(),
		Language:    repo.GetLanguage(),
		Archived:    repo.GetArchived(),
		Private:     repo.GetPrivate(),
		IsTemplate:  repo.GetIsTemplate(),
		Fork:        repo.GetFork(),
		UpdatedAt:   repo.GetUpdatedAt().Time,
	}, nil
}

func convertRepositoryFragment(host string, f githubv4.RepositoryFragment) hosting.Repository {
	parentOwner := f.GetParent().Owner
	parentName := f.GetParent().Name
	var parentRepo *hosting.ParentRepository
	if parentOwner != nil && parentName != "" {
		parentOwnerLogin := parentOwner.GetLogin()
		parentRepo = &hosting.ParentRepository{
			Ref:      repository.NewReference(host, parentOwnerLogin, parentName),
			CloneURL: convertSSHToHTTPS(f.GetParent().SshUrl),
		}
	}

	return hosting.Repository{
		Ref:         repository.NewReference(host, f.Owner.GetLogin(), f.Name),
		URL:         f.GetUrl(),
		CloneURL:    convertSSHToHTTPS(f.GetSshUrl()),
		UpdatedAt:   f.UpdatedAt,
		Parent:      parentRepo,
		Description: f.Description,
		Homepage:    f.GetHomepageUrl(),
		Language:    f.GetPrimaryLanguage().Name,
		Archived:    f.GetIsArchived(),
		Private:     f.GetIsPrivate(),
		IsTemplate:  f.GetIsTemplate(),
		Fork:        f.GetIsFork(),
	}
}

// Helper functions to convert between types

// convertOwnerAffiliations converts []hosting.RepositoryAffiliation to []githubv4.RepositoryAffiliation
func convertOwnerAffiliations(affiliations []hosting.RepositoryAffiliation) ([]githubv4.RepositoryAffiliation, error) {
	if len(affiliations) == 0 {
		return nil, nil
	}
	result := make([]githubv4.RepositoryAffiliation, len(affiliations))
	for i, affiliation := range affiliations {
		if err := typ.Remap(&(result[i]), map[hosting.RepositoryAffiliation]githubv4.RepositoryAffiliation{
			hosting.RepositoryAffiliationOwner:              githubv4.RepositoryAffiliationOwner,
			hosting.RepositoryAffiliationCollaborator:       githubv4.RepositoryAffiliationCollaborator,
			hosting.RepositoryAffiliationOrganizationMember: githubv4.RepositoryAffiliationOrganizationMember,
		}, affiliation); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func convertSSHToHTTPS(sshURL string) string {
	// Convert SSH URL to HTTPS URL
	if sshURL == "" {
		return ""
	}

	// Parse SSH URL format: git@github.com:username/repo.git
	parts := strings.Split(sshURL, "@")
	if len(parts) != 2 {
		return sshURL // Not in expected format, return as is
	}

	hostAndPath := strings.Split(parts[1], ":")
	if len(hostAndPath) != 2 {
		return sshURL // Not in expected format, return as is
	}

	host := hostAndPath[0]
	path := strings.TrimSuffix(hostAndPath[1], ".git")

	return "https://" + host + "/" + path
}

var _ hosting.HostingService = (*HostingService)(nil)
