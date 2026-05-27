package service

import (
	"context"

	"github.com/luisedmc/ghcmd/domain"
)

// RepoProvider abstracts GitHub repository operations.
type RepoProvider interface {
	Search(ctx context.Context, owner, name string) (*domain.Repository, error)
	Create(ctx context.Context, name string, private bool) (url string, err error)
}

// RepoService orchestrates repository operations.
type RepoService struct {
	repos RepoProvider
}

func NewRepoService(repos RepoProvider) *RepoService {
	return &RepoService{repos: repos}
}

// Search looks up a repository by owner and name.
func (s *RepoService) Search(ctx context.Context, owner, name string) (*domain.Repository, error) {
	return s.repos.Search(ctx, owner, name)
}

// Create creates a new repository and returns its URL.
func (s *RepoService) Create(ctx context.Context, name string, private bool) (string, error) {
	return s.repos.Create(ctx, name, private)
}
