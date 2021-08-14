package store

import (
	"database/sql"
)

type Post struct {
	ID     int
	Text   string
	UserID int
}

type SqlPostStore struct {
	db *sql.DB
}

func (p *SqlPostStore) GetAll() ([]Post, error) {
	query := "select id, text, user_id from post order by id desc limit 50"
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Text, &post.UserID)
		if err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}
