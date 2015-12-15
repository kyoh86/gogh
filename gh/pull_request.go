package gh

import "github.com/google/go-github/github"

type PullRequestSortProperty string

const (
	PullRequestSortPropertyClosed      = PullRequestSortProperty("closed")
	PullRequestSortPropertyCreated     = PullRequestSortProperty("created")
	PullRequestSortPropertyUpdated     = PullRequestSortProperty("updated")
	PullRequestSortPropertyPopularity  = PullRequestSortProperty("popularity")
	PullRequestSortPropertyLongRunning = PullRequestSortProperty("long-running")
)

type PullRequests struct {
	Order PullRequestSortProperty
	Array []github.PullRequest
}

func NewPullRequests(array []github.PullRequest) *PullRequests {
	return &PullRequests{Array: array}
}

func (p *PullRequests) Len() int {
	return len(p.Array)
}

func (p *PullRequests) Swap(i, j int) {
	p.Array[i], p.Array[j] = p.Array[j], p.Array[i]
}

func (p *PullRequests) Less(i, j int) bool {
	switch p.Order {
	case PullRequestSortPropertyClosed:
		if p.Array[i].ClosedAt == nil {
			return p.Array[j].ClosedAt != nil
		}
		if p.Array[j].ClosedAt == nil {
			return false
		}
		p.Array[i].ClosedAt.Before(*p.Array[j].ClosedAt)
	default:
		return false
		//TODO: case PullRequestSortPropertyCreated:
		//TODO: case PullRequestSortPropertyUpdated:
		//TODO: case PullRequestSortPropertyPopularity:
		//TODO: case PullRequestSortPropertyLongRunning:
	}
	return false
}
