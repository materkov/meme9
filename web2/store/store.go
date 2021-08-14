package store

import (
	"database/sql"
)

type Store struct {
	Post *SqlPostStore
	User *SqlUserStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Post: &SqlPostStore{db: db},
		User: &SqlUserStore{db: db},
	}
}
