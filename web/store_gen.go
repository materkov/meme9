package main

import "database/sql"

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Get(ids []int) ([]*Post, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "select coalesce(id, 0), coalesce(user_id, 0), coalesce(date, 0), coalesce(text, ''), coalesce(photo_id, 0) from post where id in (" + idsStr(ids) + ")"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Post, 0)
	for rows.Next() {
		obj := Post{}
		err := rows.Scan(&obj.ID, &obj.UserID, &obj.Date, &obj.Text, &obj.PhotoID)
		if err != nil {
			return nil, err
		}
		result = append(result, &obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Get(ids []int) ([]*User, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "select coalesce(id, 0), coalesce(name, ''), coalesce(avatar_id, 0), coalesce(vk_id, 0), coalesce(vk_avatar, '') from user where id in (" + idsStr(ids) + ")"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*User, 0)
	for rows.Next() {
		obj := User{}
		err := rows.Scan(&obj.ID, &obj.Name, &obj.AvatarID, &obj.VkID, &obj.VkAvatar)
		if err != nil {
			return nil, err
		}
		result = append(result, &obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type TokenStore struct {
	db *sql.DB
}

func (s *TokenStore) Get(ids []int) ([]*Token, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "select coalesce(id, 0), coalesce(token, ''), coalesce(user_id, 0) from token where id in (" + idsStr(ids) + ")"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Token, 0)
	for rows.Next() {
		obj := Token{}
		err := rows.Scan(&obj.ID, &obj.Token, &obj.UserID)
		if err != nil {
			return nil, err
		}
		result = append(result, &obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type PhotoStore struct {
	db *sql.DB
}

func (s *PhotoStore) Get(ids []int) ([]*Photo, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "select coalesce(id, 0), coalesce(user_id, 0), coalesce(path, '') from photo where id in (" + idsStr(ids) + ")"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Photo, 0)
	for rows.Next() {
		obj := Photo{}
		err := rows.Scan(&obj.ID, &obj.UserID, &obj.Path)
		if err != nil {
			return nil, err
		}
		result = append(result, &obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Get(ids []int) ([]*Comment, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "select coalesce(id, 0), coalesce(post_id, 0), coalesce(user_id, 0), coalesce(text, ''), coalesce(date, 0) from comment where id in (" + idsStr(ids) + ")"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Comment, 0)
	for rows.Next() {
		obj := Comment{}
		err := rows.Scan(&obj.ID, &obj.PostID, &obj.UserID, &obj.Text, &obj.Date)
		if err != nil {
			return nil, err
		}
		result = append(result, &obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type AllStores struct {
	Post    *PostStore
	User    *UserStore
	Token   *TokenStore
	Photo   *PhotoStore
	Comment *CommentStore
}

func NewAllStores(db *sql.DB) *AllStores {
	return &AllStores{
		Post:    &PostStore{db: db},
		User:    &UserStore{db: db},
		Token:   &TokenStore{db: db},
		Photo:   &PhotoStore{db: db},
		Comment: &CommentStore{db: db},
	}
}
