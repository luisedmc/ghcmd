package main

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// FetchToken validates a Github token and returns it with an error if invalid
func FetchToken(githubKey string) (string, error) {
	if githubKey == "" {
		return "", ErrTokenEmpty
	}

	statusCode, err := TestToken(githubKey)
	if err != nil {
		return "", ErrTokenTest
	}

	if statusCode == http.StatusUnauthorized {
		return "", ErrTokenInvalid
	}

	return githubKey, nil
}

// TestToken performs a request to the Github API to check if the token is valid
func TestToken(githubKey string) (int, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "token "+githubKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// Token returns a token source
func TokenSource(tokenInput string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenInput},
	)
}

// TokenClient returns a HTTP Client from a context and a token source
func TokenClient(ctx context.Context, ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, ts)
}
