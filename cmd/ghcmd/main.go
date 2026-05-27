package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luisedmc/ghcmd/adapter/githubapi"
	"github.com/luisedmc/ghcmd/adapter/storage"
	"github.com/luisedmc/ghcmd/service"
	"github.com/luisedmc/ghcmd/tui"
)

func main() {
	// adapters
	tokenStore, err := storage.NewTokenStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening token store: %v\n", err)
		os.Exit(1)
	}

	githubAuth := githubapi.NewGithubAuth()

	// services
	authService := service.NewAuthService(githubAuth, tokenStore)

	repoFactory := func(ctx context.Context, token string) *service.RepoService {
		ts := githubapi.TokenSource(token)
		tc := githubapi.TokenClient(ctx, ts)
		repos := githubapi.NewGithubRepos(tc)
		return service.NewRepoService(repos)
	}

	// read stored token
	storedToken, err := authService.LoadStoredToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stored token: %v\n", err)
		os.Exit(1)
	}

	m := tui.NewModel(authService, repoFactory, storedToken)
	defer m.Close()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
