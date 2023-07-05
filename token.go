package main

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"golang.org/x/oauth2"
)

// FetchToken returns a Github token, a status message and a status bool
func FetchToken(githubKey string) (string, string, bool) {
	if githubKey == "" {
		return "", "Unwritten Token", false
	}

	res := TestToken(githubKey)

	// Check if response contains "Bad credentials" == invalid API key
	if strings.Contains(string(res), "Bad credentials") {
		return "", "Invalid Token", false
	}

	return githubKey, "Valid Token", true
}

// TestToken performs a request to the Github API to check if the token is valid
func TestToken(githubKey string) []byte {
	curlCmd := exec.Command("curl", "-v", "-H", fmt.Sprintf("Authorization: token %s", githubKey), "https://api.github.com/user/issues")
	output, _ := curlCmd.CombinedOutput()
	_ = curlCmd.Run()

	return output
}

// Token returns a token source
func TokenSource(tokenInput string) (oauth2.TokenSource, error) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenInput},
	)

	return tokenSource, nil
}

// TokenClient returns a HTTP Client from a context and a token source
func TokenClient(ctx context.Context, ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, ts)
}
