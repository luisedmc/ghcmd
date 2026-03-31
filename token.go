package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/luisedmc/ghcmd/model"
	"golang.org/x/oauth2"
)

type githubUserResponse struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// FetchToken validates a Github token and returns it with the authenticated user data.
func FetchToken(githubKey string) (string, *model.User, error) {
	if githubKey == "" {
		return "", nil, ErrTokenEmpty
	}

	_, user, err := TestToken(githubKey)
	if err != nil {
		return "", nil, err
	}

	return githubKey, user, nil
}

// TestToken performs a request to the Github API to check if the token is valid
// and returns the authenticated user data on success.
func TestToken(githubKey string) (int, *model.User, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "token "+githubKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return resp.StatusCode, nil, ErrTokenInvalid
		case http.StatusForbidden:
			return resp.StatusCode, nil, ErrTokenForbidden
		case http.StatusTooManyRequests:
			return resp.StatusCode, nil, ErrTokenRateLimited
		default:
			if resp.StatusCode >= http.StatusInternalServerError {
				return resp.StatusCode, nil, ErrTokenServerError
			}
			return resp.StatusCode, nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}
	}

	var ghUser githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return resp.StatusCode, nil, err
	}

	user := &model.User{
		Login:       ghUser.Login,
		Name:        ghUser.Name,
		PublicRepos: ghUser.PublicRepos,
	}

	return resp.StatusCode, user, nil
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
