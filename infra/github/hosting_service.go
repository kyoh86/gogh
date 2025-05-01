package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

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

func getClient(ctx context.Context, host string, token *Token) *Connection {
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

// NewRemoteService creates a new RemoteService instance
func NewRemoteService(tokenService auth.TokenService) *HostingService {
	return &HostingService{
		tokenService: tokenService,
		knownOwners:  map[string]string{},
	}
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
		adaptor, err := NewAdaptor(ctx, entry.Host, &entry.Token)
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
func (s *HostingService) ListRepository(ctx context.Context, opt *hosting.ListRepositoryOptions) (list []*hosting.Repository, _ error) {
	for _, entry := range s.tokenService.Entries() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		conn := getClient(ctx, entry.Host, &entry.Token)
		var after string
		for {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			repos, err := listReposOptions(ctx, conn.gqlClient, opt, after)
			if err != nil {
				return nil, err
			}
			repositories := repos.Viewer.Repositories
			for _, edge := range repositories.Edges {
				list = append(list, util.Ptr(ingestRepositoryFragment(entry.Host, edge.Node.RepositoryFragment)))
				if opt.Limit > 0 && len(list) >= opt.Limit {
					break
				}
			}
			page := repositories.PageInfo.PageInfoFragment
			if !page.HasNextPage {
				break
			}
			after = page.EndCursor
		}
		if opt.Limit > 0 && len(list) >= opt.Limit {
			break
		}
	}
	return list, nil
}

func listReposOptions(ctx context.Context, client graphql.Client, opt *hosting.ListRepositoryOptions, after string) (*githubv4.ListReposResponse, error) {
	var limit int
	switch {
	case opt.Limit == 0:
		limit = RepoListMaxLimitPerPage
	case opt.Limit > RepoListMaxLimitPerPage:
		limit = RepoListMaxLimitPerPage
	default:
		limit = opt.Limit
	}
	return githubv4.ListRepos(
		ctx,
		client,
		limit,
		after,
		opt.IsFork.AsBoolPtr(),
		convertPrivacy(opt.Privacy),
		convertOwnerAffiliations(opt.OwnerAffiliations),
		convertRepositoryOrder(opt.OrderBy),
		convertBooleanFilter(opt.IsArchived),
	)
}

func ingestRepositoryFragment(host string, f githubv4.RepositoryFragment) hosting.Repository {
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

func (s *HostingService) ListRepositoryAsync(
	raw context.Context,
	listOption *hosting.ListRepositoryOptions,
) (<-chan *hosting.Repository, <-chan error) {
	repoChan := make(chan *hosting.Repository)
	errChan := make(chan error)
	raise := func(err error) {
		if !errors.Is(err, context.Canceled) {
			errChan <- err
		}
	}
	ctx, cancel := context.WithCancel(raw)
	defer cancel()

	var count int32
	var wg sync.WaitGroup
	for _, entry := range s.tokenService.Entries() {
		wg.Add(1)
		go func(entry auth.TokenEntry) {
			defer wg.Done()

			conn := getClient(ctx, entry.Host, &entry.Token)
			var after string
			for {
				if err := ctx.Err(); err != nil {
					raise(err)
					return
				}
				repos, err := listReposOptions(ctx, conn.gqlClient, listOption, after)
				if err != nil {
					raise(err)
					return
				}
				repositories := repos.Viewer.Repositories
				for _, edge := range repositories.Edges {
					select {
					case repoChan <- util.Ptr(ingestRepositoryFragment(entry.Host, edge.Node.RepositoryFragment)):
						if int(atomic.AddInt32(&count, 1)) >= listOption.Limit {
							cancel()
							return
						}
					case <-ctx.Done():
						if err := ctx.Err(); err != nil {
							raise(err)
						}
						return
					}
				}
				page := repositories.PageInfo.PageInfoFragment
				if !page.HasNextPage {
					break
				}
				after = page.EndCursor
			}
		}(entry)
	}

	go func() {
		wg.Wait()
		close(repoChan)
		close(errChan)
	}()

	return repoChan, errChan
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
