package db

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// DBConnect returns the singleton instance of the database
func DBConnect(dbURL string) (*sql.DB, error) {
	if dbURL == "" {
		return nil, errors.New("dbURL is empty")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// LoadSqlTestFile read the sql test file and exec it
func LoadSqlTestFile(t *testing.T, tx *sql.DB, sqlFile string) {
	b, err := ioutil.ReadFile(sqlFile)
	require.NoError(t, err)

	_, err = tx.Exec(string(b))
	require.NoError(t, err)
}
