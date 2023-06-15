package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func Token() (oauth2.TokenSource, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ghKey := os.Getenv("GITHUB_API_KEY")

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghKey},
	)

	return tokenSource, nil
}

func TokenClient(ctx context.Context, ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, ts)
}
