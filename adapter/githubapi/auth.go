package githubapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/luisedmc/ghcmd/domain"
)

type githubUserResponse struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	PublicRepos int    `json:"public_repos"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GithubAuth implements service.AuthProvider using the GitHub REST API.
type GithubAuth struct{}

func NewGithubAuth() *GithubAuth {
	return &GithubAuth{}
}

// ValidateToken checks whether a GitHub token is valid and returns the
// authenticated user on success.
func (g *GithubAuth) ValidateToken(token string) (*domain.User, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, domain.ErrTokenInvalid
		case http.StatusForbidden:
			return nil, domain.ErrTokenForbidden
		case http.StatusTooManyRequests:
			return nil, domain.ErrTokenRateLimited
		default:
			if resp.StatusCode >= http.StatusInternalServerError {
				return nil, domain.ErrTokenServerError
			}
			return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}
	}

	var ghUser githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, err
	}

	return &domain.User{
		Login:       ghUser.Login,
		Name:        ghUser.Name,
		PublicRepos: ghUser.PublicRepos,
	}, nil
}

// TokenSource returns an OAuth2 static token source.
func TokenSource(token string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
}

// TokenClient returns a HTTP Client from a context and a token source.
func TokenClient(ctx context.Context, ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, ts)
}
