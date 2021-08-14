package store

import (
	"database/sql"
)

type Store struct {
	Post *SqlPostStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Post: &SqlPostStore{
			db: db,
		},
	}
}
