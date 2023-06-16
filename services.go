package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/go-github/v53/github"
)

func GithubClient(tokenClient *http.Client) *github.Client {
	githubClient := github.NewClient(tokenClient)

	return githubClient
}

func SearchRepository(githubClient *github.Client, user string, repositoryName string) *Repository {
	ctx := context.Background()

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
