package app

import (
	"database/sql"
	"testing"

	"github.com/materkov/meme9/web/store"
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
	ObjectStore = store.NewObjectStore(db)
}
