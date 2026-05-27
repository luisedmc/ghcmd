package githubapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v84/github"

	"github.com/luisedmc/ghcmd/domain"
)

// GithubRepos implements service.RepoProvider using the GitHub API client.
type GithubRepos struct {
	client *github.Client
}

func NewGithubRepos(tokenClient *http.Client) *GithubRepos {
	return &GithubRepos{client: github.NewClient(tokenClient)}
}

// Search looks up a specific repository by owner and name.
func (r *GithubRepos) Search(ctx context.Context, owner, name string) (*domain.Repository, error) {
	repository, _, err := r.client.Repositories.Get(ctx, owner, name)
	if err != nil {
		var ghErr *github.ErrorResponse
		if errors.As(err, &ghErr) && ghErr.Response.StatusCode == http.StatusNotFound {
			return nil, domain.ErrSearchNotFound
		}
		return nil, fmt.Errorf("failed to search repository: %w", err)
	}

	if repository.GetPrivate() {
		return nil, domain.ErrSearchNotFound
	}

	return &domain.Repository{
		Name:        repository.GetName(),
		Owner:       repository.GetOwner().GetLogin(),
		OwnerURL:    repository.GetOwner().GetHTMLURL(),
		Description: repository.GetDescription(),
		URL:         repository.GetHTMLURL(),
		Stars:       repository.GetStargazersCount(),
		Forks:       repository.GetForksCount(),
		Language:    repository.GetLanguage(),
		OpenIssues:  repository.GetOpenIssuesCount(),
		CreatedAt:   repository.GetCreatedAt().Format("Jan 2, 2006"),
		License:     repository.GetLicense().GetName(),
	}, nil
}

// Create creates a new repository in the authenticated user's account.
func (r *GithubRepos) Create(ctx context.Context, name string, private bool) (string, error) {
	newRepository := &github.Repository{
		Name:    &name,
		Private: &private,
	}

	res, _, err := r.client.Repositories.Create(ctx, "", newRepository)
	if err != nil {
		var ghErr *github.ErrorResponse
		if errors.As(err, &ghErr) {
			switch ghErr.Response.StatusCode {
			case http.StatusUnprocessableEntity:
				return "", domain.ErrRepoAlreadyExists
			case http.StatusNotFound:
				return "", domain.ErrRepoUnauthorized
			}
		}
		return "", domain.ErrRepoCreateFailed
	}

	return res.GetHTMLURL(), nil
}
