package github

import (
	"context"
	"fmt"
	"iter"
	"net/url"
	"strings"

	"github.com/Khan/genqlient/graphql"
	github "github.com/google/go-github/v69/github"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/infra/githubv4"
	"github.com/kyoh86/gogh/v3/util"
	"golang.org/x/oauth2"
)

type HostingService struct {
	// GitHub client, etc.
	tokenService auth.TokenService
	knownOwners  map[string]string
}

type Connection struct {
	tokenSource oauth2.TokenSource
	restClient  *github.Client
	gqlClient   graphql.Client
}

const (
	GlobalHost    = "github.com"
	GlobalAPIHost = "api.github.com"
)

const ClientID = "Ov23li6aEWIxek6F8P5L"

func getClient(ctx context.Context, host string, token *auth.Token) *Connection {
	var source oauth2.TokenSource
	if token != nil {
		source = oauth2.ReuseTokenSource(token, &tokenSource{ctx: ctx, host: host, token: token})
	}
	httpClient := oauth2.NewClient(ctx, source)
	if host == GlobalHost || host == GlobalAPIHost {
		return &Connection{
			tokenSource: source,
			restClient:  github.NewClient(httpClient),
			gqlClient:   graphql.NewClient("https://"+GlobalAPIHost+"/graphql", httpClient),
		}
	}
	baseRESTURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/v3",
	}
	uploadRESTURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/uploads",
	}
	baseGQLURL := &url.URL{
		Scheme: "https://",
		Host:   host,
		Path:   "/api/graphql",
	}
	restClient, err := github.NewClient(httpClient).WithEnterpriseURLs(baseRESTURL.String(), uploadRESTURL.String())
	if err != nil {
		// NOTE: WithEnterpriseURLs returns error if the URLs are not valid.
		// We assume that the URLs are valid.
		panic(err)
	}
	return &Connection{
		tokenSource: source,
		restClient:  restClient,
		gqlClient:   graphql.NewClient(baseGQLURL.String(), httpClient),
	}
}

// NewHostingService creates a new HostingService instance
func NewHostingService(tokenService auth.TokenService) *HostingService {
	return &HostingService{
		tokenService: tokenService,
		knownOwners:  map[string]string{},
	}
}

// GetOauth2Config implements hosting.HostingService.
func (s *HostingService) GetOauth2Config(ctx context.Context, host string) (*oauth2.Config, error) {
	return &oauth2.Config{
		ClientID: ClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:       fmt.Sprintf("https://%s/login/oauth/authorize", host),
			TokenURL:      fmt.Sprintf("https://%s/login/oauth/access_token", host),
			DeviceAuthURL: fmt.Sprintf("https://%s/login/device/code", host),
		},
		Scopes: []string{string(github.ScopeRepo), string(github.ScopeDeleteRepo)},
	}, nil
}

// GetAuthenticatedUserName implements hosting.HostingService.
func (s *HostingService) GetAuthenticatedUserName(ctx context.Context, host string, token *oauth2.Token) (string, error) {
	conn := getClient(ctx, host, token)
	user, _, err := conn.restClient.Users.Get(ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get authenticated user: %w", err)
	}
	return user.GetLogin(), nil
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
	words := strings.SplitN(u.Path, "/", 2)
	if len(words) < 2 {
		return nil, fmt.Errorf("invalid path: %q", u.Path)
	}
	return util.Ptr(repository.NewReference(u.Host, words[0], strings.TrimSuffix(words[1], ".git"))), nil
}

// GetTokenFor cache requested token for the host and owner
func (s *HostingService) GetTokenFor(ctx context.Context, reference repository.Reference) (string, auth.Token, error) {
	key := strings.Join([]string{reference.Host(), reference.Owner()}, "/")
	tokenOwner, ok := s.knownOwners[key]
	if ok {
		_, token, err := s.getTokenForCore(ctx, tokenOwner, reference.Name())
		return tokenOwner, token, err
	}
	tokenOwner, token, err := s.getTokenForCore(ctx, reference.Host(), reference.Owner())
	if err == nil {
		s.knownOwners[key] = tokenOwner
	}
	return tokenOwner, token, err
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
		adaptor, err := NewAdaptor(ctx, entry.Host, util.Ptr((Token)(entry.Token)))
		if err != nil {
			continue // Try next token if this one fails
		}

		// Check if this user is a member of the target organization
		ok, err := adaptor.MemberOf(ctx, owner)
		if err != nil {
			continue // Try next token if membership check fails
		}

		if ok {
			// Found a working token
			return entry.Owner, entry.Token, nil
		}
	}

	return "", auth.Token{}, fmt.Errorf("no valid token found for %s/%s", host, owner)
}

// GetRepository retrieves repository information from a remote source
func (s *HostingService) GetRepository(ctx context.Context, reference repository.Reference) (*hosting.Repository, error) {
	_, token, err := s.GetTokenFor(ctx, reference)
	if err != nil {
		return nil, fmt.Errorf("failed to get token for %s/%s: %w", reference.Host(), reference.Owner(), err)
	}
	conn := getClient(ctx, reference.Host(), &token)
	ghRepo, _, err := conn.restClient.Repositories.Get(ctx, reference.Owner(), reference.Name())
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
			return nil, fmt.Errorf("failed to parse parent HTML URL: %w", err)
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

			conn := getClient(ctx, entry.Host, &entry.Token)
			var after string

			for {
				if err := ctx.Err(); err != nil {
					yield(nil, err)
					return
				}

				repos, err := githubv4.ListRepos(
					ctx,
					conn.gqlClient,
					limit,
					after,
					opts.IsFork.AsBoolPtr(),
					convertPrivacy(opts.Privacy),
					convertOwnerAffiliations(opts.OwnerAffiliations),
					convertRepositoryOrder(opts.OrderBy),
					convertBooleanFilter(opts.IsArchived),
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
					if !yield(util.Ptr(convertRepositoryFragment(entry.Host, edge.Node.RepositoryFragment)), nil) {
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

func (s *HostingService) CreateRepository(
	ctx context.Context,
	ref repository.Reference,
	opts hosting.CreateRepositoryOptions,
) (*hosting.Repository, error) {
	user, token, err := s.GetTokenFor(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get token for %s/%s: %w", ref.Host(), ref.Owner(), err)
	}
	conn := getClient(ctx, ref.Host(), &token)
	org := ""
	if user != ref.Owner() {
		org = ref.Owner()
	}
	repo, _, err := conn.restClient.Repositories.Create(ctx, org, &github.Repository{
		Name:                util.NilablePtr(ref.Name()),
		Description:         util.NilablePtr(opts.Description),
		Homepage:            util.NilablePtr(opts.Homepage),
		Private:             util.NilablePtr(opts.Private),
		HasIssues:           util.FalsePtr(opts.DisableIssues),
		HasProjects:         util.FalsePtr(opts.DisableProjects),
		HasWiki:             util.FalsePtr(opts.DisableWiki),
		HasDownloads:        util.FalsePtr(opts.DisableDownloads),
		IsTemplate:          util.NilablePtr(opts.IsTemplate),
		TeamID:              util.NilablePtr(opts.TeamID),
		AutoInit:            util.NilablePtr(opts.AutoInit),
		GitignoreTemplate:   util.NilablePtr(opts.GitignoreTemplate),
		LicenseTemplate:     util.NilablePtr(opts.LicenseTemplate),
		AllowSquashMerge:    util.FalsePtr(opts.PreventSquashMerge),
		AllowMergeCommit:    util.FalsePtr(opts.PreventMergeCommit),
		AllowRebaseMerge:    util.FalsePtr(opts.PreventRebaseMerge),
		DeleteBranchOnMerge: util.NilablePtr(opts.DeleteBranchOnMerge),
	})
	if err != nil {
		return nil, err
	}
	return convertRepository(ref, repo)
}

func (s *HostingService) CreateRepositoryFromTemplate(
	ctx context.Context,
	ref repository.Reference,
	template repository.Reference,
	opts hosting.CreateRepositoryFromTemplateOptions,
) (*hosting.Repository, error) {
	user, token, err := s.GetTokenFor(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get token for %s/%s: %w", ref.Host(), ref.Owner(), err)
	}
	conn := getClient(ctx, ref.Host(), &token)
	req := github.TemplateRepoRequest{
		Name:               util.Ptr(ref.Name()),
		Description:        &opts.Description,
		IncludeAllBranches: &opts.IncludeAllBranches,
		Private:            &opts.Private,
	}
	if user != ref.Owner() {
		req.Owner = util.Ptr(ref.Owner())
	}
	repo, _, err := conn.restClient.Repositories.CreateFromTemplate(
		ctx,
		template.Owner(),
		template.Name(),
		&req,
	)
	if err != nil {
		return nil, err
	}
	return convertRepository(ref, repo)
}

// DeleteRepository deletes a repository from a remote source
func (s *HostingService) DeleteRepository(ctx context.Context, reference repository.Reference) error {
	_, token, err := s.GetTokenFor(ctx, reference)
	if err != nil {
		return fmt.Errorf("failed to get token for %s/%s: %w", reference.Host(), reference.Owner(), err)
	}
	conn := getClient(ctx, reference.Host(), &token)
	_, err = conn.restClient.Repositories.Delete(ctx, reference.Owner(), reference.Name())
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
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
	user, token, err := s.GetTokenFor(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to get token for %s/%s: %w", ref.Host(), ref.Owner(), err)
	}
	conn := getClient(ctx, ref.Host(), &token)
	ghOpts := &github.RepositoryCreateForkOptions{
		Name:              target.Name(),
		DefaultBranchOnly: opts.DefaultBranchOnly,
	}
	if user != target.Owner() {
		ghOpts.Organization = target.Owner()
	}
	fork, _, err := conn.restClient.Repositories.CreateFork(ctx, ref.Owner(), ref.Name(), ghOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create fork: %w", err)
	}
	return convertRepository(ref, fork)
}

func convertRepository(ref repository.Reference, repo *github.Repository) (*hosting.Repository, error) {
	var parent *hosting.ParentRepository
	if raw := repo.GetParent(); raw != nil {
		u, err := url.Parse(raw.GetHTMLURL())
		if err != nil {
			return nil, err
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
		Ref:         ref,
		URL:         *repo.HTMLURL,
		Parent:      parent,
		Description: repo.GetDescription(),
		Homepage:    repo.GetHomepage(),
		Language:    repo.GetLanguage(),
		Archived:    repo.GetArchived(),
		Private:     repo.GetPrivate(),
		IsTemplate:  repo.GetIsTemplate(),
		Fork:        repo.GetFork(),
	}, nil
}
func convertRepositoryFragment(host string, f githubv4.RepositoryFragment) hosting.Repository {
	parentOwner := f.GetParent().Owner.GetLogin()
	parentName := f.GetParent().Name
	var parentRepo *hosting.ParentRepository
	if parentOwner != "" && parentName != "" {
		parentRepo = &hosting.ParentRepository{
			Ref:      repository.NewReference(host, parentOwner, parentName),
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

// convertPrivacy converts hosting.RepositoryPrivacy to githubv4.RepositoryPrivacy
func convertPrivacy(privacy hosting.RepositoryPrivacy) githubv4.RepositoryPrivacy {
	switch privacy {
	case hosting.RepositoryPrivacyPublic:
		return githubv4.RepositoryPrivacyPublic
	case hosting.RepositoryPrivacyPrivate:
		return githubv4.RepositoryPrivacyPrivate
	default:
		return githubv4.RepositoryPrivacy("")
	}
}

// convertOwnerAffiliations converts []hosting.RepositoryAffiliation to []githubv4.RepositoryAffiliation
func convertOwnerAffiliations(affiliations []hosting.RepositoryAffiliation) []githubv4.RepositoryAffiliation {
	if len(affiliations) == 0 {
		return nil
	}

	result := make([]githubv4.RepositoryAffiliation, len(affiliations))
	for i, affiliation := range affiliations {
		switch affiliation {
		case hosting.RepositoryAffiliationOwner:
			result[i] = githubv4.RepositoryAffiliationOwner
		case hosting.RepositoryAffiliationCollaborator:
			result[i] = githubv4.RepositoryAffiliationCollaborator
		case hosting.RepositoryAffiliationOrganizationMember:
			result[i] = githubv4.RepositoryAffiliationOrganizationMember
		}
	}
	return result
}

// convertBooleanFilter converts hosting.BooleanFilter to *bool
func convertBooleanFilter(filter hosting.BooleanFilter) *bool {
	switch filter {
	case hosting.BooleanFilterTrue:
		value := true
		return &value
	case hosting.BooleanFilterFalse:
		value := false
		return &value
	default: // BooleanFilterNone
		return nil
	}
}

func convertRepositoryOrder(order hosting.RepositoryOrder) githubv4.RepositoryOrder {
	// Map the field
	var field githubv4.RepositoryOrderField
	switch order.Field {
	case hosting.RepositoryOrderFieldCreatedAt:
		field = githubv4.RepositoryOrderFieldCreatedAt
	case hosting.RepositoryOrderFieldUpdatedAt:
		field = githubv4.RepositoryOrderFieldUpdatedAt
	case hosting.RepositoryOrderFieldPushedAt:
		field = githubv4.RepositoryOrderFieldPushedAt
	case hosting.RepositoryOrderFieldName:
		field = githubv4.RepositoryOrderFieldName
	case hosting.RepositoryOrderFieldStargazers:
		field = githubv4.RepositoryOrderFieldStargazers
	default:
		// Default to created time if not recognized
		field = githubv4.RepositoryOrderFieldCreatedAt
	}

	// Map the direction
	var direction githubv4.OrderDirection
	switch order.Direction {
	case hosting.OrderDirectionAsc:
		direction = githubv4.OrderDirectionAsc
	case hosting.OrderDirectionDesc:
		direction = githubv4.OrderDirectionDesc
	default:
		// Default to descending if not recognized
		direction = githubv4.OrderDirectionDesc
	}

	return githubv4.RepositoryOrder{
		Field:     field,
		Direction: direction,
	}
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

// Ensure RemoteService implements repository.RemoteRepositoryService
var _ hosting.HostingService = (*HostingService)(nil)
