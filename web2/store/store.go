package store

import (
	"database/sql"
)

type Store struct {
	Post   *SqlPostStore
	User   *SqlUserStore
	Object *SqlObjectStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Post:   &SqlPostStore{db: db},
		User:   &SqlUserStore{db: db},
		Object: &SqlObjectStore{db: db},
	}
}
