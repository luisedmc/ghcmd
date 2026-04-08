package storage

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/luisedmc/ghcmd/service"
)

var _ service.TokenStore = (*TokenStore)(nil)

type TokenStore struct {
	path string
}

// NewTokenStore creates the config directory if needed and returns a store
// that reads/writes the token from a plain file.
func NewTokenStore() (*TokenStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(configDir, "ghcmd")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &TokenStore{path: filepath.Join(dir, "token")}, nil
}

// ReadToken reads the stored GitHub token. Returns "" if no token exists.
func (s *TokenStore) ReadToken() (string, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// SaveToken writes the GitHub token to disk with owner-only permissions.
func (s *TokenStore) SaveToken(token string) error {
	return os.WriteFile(s.path, []byte(token), 0600)
}

// Close is a no-op kept for interface consistency.
func (s *TokenStore) Close() error {
	return nil
}
