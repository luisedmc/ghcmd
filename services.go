package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v53/github"
)

// Repository represents a Github repository
type Repository struct {
	FullName    string
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
		log.Printf("Problem in getting repository information %v\n", err)
	}

	repositoryData := &Repository{
		FullName:    *repository.Owner.HTMLURL,
		Description: *repository.Description,
		URL:         *repository.HTMLURL,
	}

	return repositoryData
}

func CreateRepository(ctx context.Context, repoName string, isPrivate bool, githubClient *github.Client) {
	newRepository := &github.Repository{
		Name:    github.String(repoName),
		Private: github.Bool(isPrivate),
	}

	_, _, err := githubClient.Repositories.Create(ctx, "", newRepository)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Repository created successfully!")
	}
}
