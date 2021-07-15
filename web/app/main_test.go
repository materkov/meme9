package app

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	db.SetMaxOpenConns(1)

	_, err = db.Exec(migrations)
	require.NoError(t, err)

	// TODO replace this global
	store = NewAllStores(db)
}
