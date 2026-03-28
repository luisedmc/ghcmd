package db

import (
	"os"
	"path/filepath"
	"strings"
)

type Database struct {
	path string
}

// OpenDB creates the config directory if needed and returns a Database
// that reads/writes the token from a plain file.
func OpenDB() (*Database, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(configDir, "ghcmd")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Database{path: filepath.Join(dir, "token")}, nil
}

// GetToken reads the stored GitHub token. Returns "" if no token exists.
func (d *Database) GetToken() (string, error) {
	data, err := os.ReadFile(d.path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// SetToken writes the GitHub token to disk with owner-only permissions.
func (d *Database) SetToken(token string) error {
	return os.WriteFile(d.path, []byte(token), 0600)
}

// Close is a no-op kept for interface consistency.
func (d *Database) Close() error {
	return nil
}
