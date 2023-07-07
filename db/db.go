package db

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type Database struct {
	Conn *leveldb.DB
}

func createDB(dbPath string) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return err
	}
	// log.Println("Database created successfully.")
	defer db.Close()

	return nil
}

// OpenDB opens a database file
func OpenDB() (*Database, error) {
	_, err := os.Stat("./db/data")
	if err != nil {
		if os.IsNotExist(err) {
			err = createDB("./db/data")
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	db, err := leveldb.OpenFile("./db/data", nil)
	if err != nil {
		return nil, err
	}

	d := &Database{db}
	// log.Println("Database opened successfully.")

	return d, nil
}

// PutToken stores the token in the database
func (d *Database) PutToken(db *leveldb.DB, token string) error {
	err := db.Put([]byte("gh_token"), []byte(token), nil)
	if err != nil {
		return err
	}

	return nil
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
