package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func apiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	githubKey := os.Getenv("GITHUB_API_KEY")

	return githubKey
}

func Token() (oauth2.TokenSource, error) {
	githubKey := apiKey()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubKey},
	)

	return tokenSource, nil
}

func TokenClient(ts oauth2.TokenSource) *http.Client {
	ctx := context.Background()

	return oauth2.NewClient(ctx, ts)
}
