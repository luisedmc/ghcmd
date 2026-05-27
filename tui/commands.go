package tui

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/luisedmc/ghcmd/domain"
	"github.com/luisedmc/ghcmd/service"
)

func fetchUserCmd(auth *service.AuthService, token string) tea.Cmd {
	return func() tea.Msg {
		user, err := auth.ValidateToken(token)
		return userFetchedMsg{user: user, token: token, err: err}
	}
}

func searchRepoCmd(repos *service.RepoService, ctx context.Context, owner, name string) tea.Cmd {
	return func() tea.Msg {
		data, err := repos.Search(ctx, owner, name)
		msg := serviceResultMsg{responseData: data}
		if err != nil {
			msg.err = err
			switch {
			case errors.Is(err, domain.ErrSearchNotFound):
				msg.message = "Repository not found!"
			default:
				msg.message = "Failed to search repository!"
			}
		}
		return msg
	}
}

func createRepoCmd(repos *service.RepoService, ctx context.Context, repoName string, private bool) tea.Cmd {
	return func() tea.Msg {
		url, err := repos.Create(ctx, repoName, private)
		msg := serviceResultMsg{}
		if err != nil {
			msg.err = err
			switch {
			case errors.Is(err, domain.ErrRepoAlreadyExists):
				msg.message = "Repository already exists!"
			case errors.Is(err, domain.ErrRepoUnauthorized):
				msg.message = "Token lacks permission to create repositories!"
			case errors.Is(err, domain.ErrRepoCreateFailed):
				msg.message = "Repository creation failed!"
			}
		} else {
			msg.url = &url
			msg.message = "Repository created successfully!"
		}
		return msg
	}
}
