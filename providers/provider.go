package providers

import (
	"github.com/google/go-github/v22/github"
	gittypes "github.com/mohuk/code-review/types"
)

// Provider interface
type Provider interface {
	InitializeClient() (*github.Client, error)
	CreateReviewPullRequest(base string, head string) (*github.PullRequest, error)
	CreateReviewBranch(base, name string) (*gittypes.GitBranch, error)
	BranchExists(branch string) (bool, error)
}

// ProviderConfig struct
type ProviderConfig struct {
	Name        string
	AccessToken string
	UserAccount string
	Repository  string
}

// NewProvider is what
func NewProvider(cfg ProviderConfig) Provider {
	switch cfg.Name {
	default:
		{
			return &GithubReview{
				AccessToken: cfg.AccessToken,
				UserAccount: cfg.UserAccount,
				Repository:  cfg.Repository,
			}
		}
	}
}
