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

func apiKey() (string, string, bool) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return "", "Token Not Found", false
	}

	githubKey := os.Getenv("GITHUB_API_KEY")

	ghkey, msg, status := FetchToken(githubKey)

	return ghkey, msg, status
}

func FetchToken(githubKey string) (string, string, bool) {
	if githubKey == "" {
		return "", "Token Not Found", false
	}

	curlCmd := exec.Command("curl", "-v", "-H", fmt.Sprintf("Authorization: token %s", githubKey), "https://api.github.com/user/issues")
	output, _ := curlCmd.CombinedOutput()
	_ = curlCmd.Run()

	// Check if response contains "Bad credentials" == invalid API key
	if strings.Contains(string(output), "Bad credentials") {
		return "", "Invalid Token", false
	}

	return githubKey, "Valid Token", true
}

// Token returns a token source
func Token() (oauth2.TokenSource, error) {
	githubKey, _, _ := apiKey()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubKey},
	)

	return tokenSource, nil
}

// TokenClient returns a HTTP Client from a context and a token source
func TokenClient(ctx context.Context, ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, ts)
}
