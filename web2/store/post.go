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

func (s *SqlPostStore) GetAll() ([]*Post, error) {
	query := "select id from object where object_type = 1 order by id desc limit 50"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Text, &post.UserID)
		if err != nil {
			return nil, err
		}

		result = append(result, &post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *SqlPostStore) GetById(id int) (*Post, error) {
	p := Post{}
	err := s.db.QueryRow("select id, user_id, text from post where id = ?", id).Scan(&p.ID, &p.UserID, &p.Text)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &p, err
}

func (s *SqlPostStore) GetByUser(userID int, limit int) ([]*Post, error) {
	query := "select id, text, user_id from post where user_id = ? order by id desc limit ?"
	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Text, &post.UserID)
		if err != nil {
			return nil, err
		}

		result = append(result, &post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *SqlPostStore) Add(post *Post) error {
	result, err := s.db.Exec("insert into post(text, user_id) values (?, ?)", post.Text, post.UserID)
	if err != nil {
		return err
	}

	postID, _ := result.LastInsertId()
	post.ID = int(postID)

	return nil
}
