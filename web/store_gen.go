package main

import (
	"database/sql"
)

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

func (s *PostStore) Add(obj *Post) error {
	query := "insert into post(id, user_id, date, text, photo_id) values (?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, sql.NullInt32{Int32: int32(obj.ID), Valid: obj.ID != 0}, sql.NullInt32{Int32: int32(obj.UserID), Valid: obj.UserID != 0}, sql.NullInt32{Int32: int32(obj.Date), Valid: obj.Date != 0}, sql.NullString{String: obj.Text, Valid: obj.Text != ""}, sql.NullInt32{Int32: int32(obj.PhotoID), Valid: obj.PhotoID != 0})
	return err
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

func (s *UserStore) Add(obj *User) error {
	query := "insert into user(id, name, avatar_id, vk_id, vk_avatar) values (?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, sql.NullInt32{Int32: int32(obj.ID), Valid: obj.ID != 0}, sql.NullString{String: obj.Name, Valid: obj.Name != ""}, sql.NullInt32{Int32: int32(obj.AvatarID), Valid: obj.AvatarID != 0}, sql.NullInt32{Int32: int32(obj.VkID), Valid: obj.VkID != 0}, sql.NullString{String: obj.VkAvatar, Valid: obj.VkAvatar != ""})
	return err
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

func (s *TokenStore) Add(obj *Token) error {
	query := "insert into token(id, token, user_id) values (?, ?, ?)"
	_, err := s.db.Exec(query, sql.NullInt32{Int32: int32(obj.ID), Valid: obj.ID != 0}, sql.NullString{String: obj.Token, Valid: obj.Token != ""}, sql.NullInt32{Int32: int32(obj.UserID), Valid: obj.UserID != 0})
	return err
}

type APILogStore struct {
	db *sql.DB
}

func (s *APILogStore) Get(ids []int) ([]*APILog, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "select coalesce(id, 0), coalesce(user_id, 0), coalesce(method, ''), coalesce(request, ''), coalesce(response, '') from apilog where id in (" + idsStr(ids) + ")"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*APILog, 0)
	for rows.Next() {
		obj := APILog{}
		err := rows.Scan(&obj.ID, &obj.UserID, &obj.Method, &obj.Request, &obj.Response)
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

func (s *APILogStore) Add(obj *APILog) error {
	query := "insert into apilog(id, user_id, method, request, response) values (?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, sql.NullInt32{Int32: int32(obj.ID), Valid: obj.ID != 0}, sql.NullInt32{Int32: int32(obj.UserID), Valid: obj.UserID != 0}, sql.NullString{String: obj.Method, Valid: obj.Method != ""}, sql.NullString{String: obj.Request, Valid: obj.Request != ""}, sql.NullString{String: obj.Response, Valid: obj.Response != ""})
	return err
}

type Store struct {
	db        *sql.DB
	Post      *PostStore
	User      *UserStore
	Token     *TokenStore
	//Photo     *PhotoStore
	//Likes     *LikesStore
	//Comment   *CommentStore
	APILog    *APILogStore
	//Followers *FollowersStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:        db,
		Post:      &PostStore{db: db},
		User:      &UserStore{db: db},
		Token:     &TokenStore{db: db},
		//Photo:     &PhotoStore{db: db},
		//Likes:     &LikesStore{db: db},
		//Comment:   &CommentStore{db: db},
		APILog:    &APILogStore{db: db},
		//Followers: &FollowersStore{db: db},
	}
}
