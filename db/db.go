package db

import (
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

type Database struct {
	Conn *leveldb.DB
}

// OpenDB opens (or creates) the database at the user's config directory
func OpenDB() (*Database, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(configDir, "ghcmd", "data")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// GetToken retrieves the token from the database
func (d *Database) GetToken(db *leveldb.DB) (string, error) {
	token, err := db.Get([]byte("gh_token"), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return "", nil
		}
		return "", err
	}

	// log.Println("token: ", string(token))
	return string(token), nil
}
