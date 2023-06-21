package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func apiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return ""
	}

	githubKey := os.Getenv("GITHUB_API_KEY")

	// Github API Key curl test
	curlCmd := exec.Command("curl", "-v", "-H", fmt.Sprintf("Authorization: token %s", githubKey), "https://api.github.com/user/issues")
	output, _ := curlCmd.CombinedOutput()
	_ = curlCmd.Run()

	// Check if response contains "Bad credentials" == invalid API key
	if strings.Contains(string(output), "Bad credentials") {
		return ""
	}

	return githubKey
}

func Token() (oauth2.TokenSource, error) {
	githubKey := apiKey()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubKey},
	)

	return tokenSource, nil
}

func TokenClient(ctx context.Context, ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, ts)
}
