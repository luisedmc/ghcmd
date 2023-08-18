package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/go-github/v53/github"
)

// Repository represents a Github repository
type Repository struct {
	Owner       string
	OwnerURL    string
	Description string
	URL         string
}

// GithubClient returns a new Github client
func GithubClient(tokenClient *http.Client) *github.Client {
	githubClient := github.NewClient(tokenClient)

	return githubClient
}

// SearchRepository performs a search for a specific repository from an user and returns the repository information
func SearchRepository(ctx context.Context, githubClient *github.Client, user string, repositoryName string) *Repository {
	repository, _, err := githubClient.Repositories.Get(ctx, user, repositoryName)
	if err != nil {
		return nil
	}

	if *repository.Private {
		return nil
	}

	checkRepositoryInfoNil(repository)

	repositoryData := &Repository{
		Owner:       *repository.Owner.Login,
		OwnerURL:    *repository.Owner.HTMLURL,
		Description: *repository.Description,
		URL:         *repository.HTMLURL,
	}

	return repositoryData
}

func checkRepositoryInfoNil(repository *github.Repository) {
	if repository.Description == nil {
		repository.Description = github.String("No description provided.")
	}

	if repository.HTMLURL == nil || *repository.HTMLURL == "" {
		repository.HTMLURL = github.String("No URL provided.")
	}
}

// CreateRepository creates a new repository in the user account.
func CreateRepository(ctx context.Context, githubClient *github.Client, repoName string, isPrivate string) (*string, string, error) {
	isPrivateBool := false
	if isPrivate == "y" {
		isPrivateBool = true
	}

	newRepository := &github.Repository{
		Name:    github.String(repoName),
		Private: github.Bool(isPrivateBool),
	}

	res, _, err := githubClient.Repositories.Create(ctx, "", newRepository)
	if err != nil {
		if strings.Contains(err.Error(), "422") {
			return nil, "Repository already exists!", err
		}
		return nil, "Repository creation failed!", err
	}

	return res.HTMLURL, "Repository created successfully!", nil
}
