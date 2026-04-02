package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v84/github"

	"github.com/luisedmc/ghcmd/model"
)

// GithubClient returns a new Github client
func GithubClient(tokenClient *http.Client) *github.Client {
	return github.NewClient(tokenClient)
}

// SearchRepository performs a search for a specific repository from an user and returns the repository information
func SearchRepository(ctx context.Context, githubClient *github.Client, user string, repositoryName string) (*model.Repository, error) {
	repository, _, err := githubClient.Repositories.Get(ctx, user, repositoryName)
	if err != nil {
		var ghErr *github.ErrorResponse
		if errors.As(err, &ghErr) && ghErr.Response.StatusCode == http.StatusNotFound {
			return nil, ErrSearchNotFound
		}
		return nil, fmt.Errorf("failed to search repository: %w", err)
	}

	if repository.GetPrivate() {
		return nil, ErrSearchNotFound
	}

	repositoryData := &model.Repository{
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
	}

	return repositoryData, nil
}

// CreateRepository creates a new repository in the user account.
func CreateRepository(ctx context.Context, githubClient *github.Client, repoName string, isPrivate string) (*string, error) {
	isPrivateBool := false
	if isPrivate == "y" {
		isPrivateBool = true
	} else if isPrivate != "n" && isPrivate != "" {
		return nil, ErrInvalidPrivateInput
	}

	newRepository := &github.Repository{
		Name:    &repoName,
		Private: &isPrivateBool,
	}

	res, _, err := githubClient.Repositories.Create(ctx, "", newRepository)
	if err != nil {
		var ghErr *github.ErrorResponse
		if errors.As(err, &ghErr) {
			switch ghErr.Response.StatusCode {
			case http.StatusUnprocessableEntity:
				return nil, ErrRepoAlreadyExists
			case http.StatusNotFound:
				return nil, ErrRepoUnauthorized
			}
		}
		return nil, ErrRepoCreateFailed
	}

	return res.HTMLURL, nil
}
