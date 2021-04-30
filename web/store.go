package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Post struct {
	ID      int    `db:"id"`
	UserID  int    `db:"user_id"`
	Date    int    `db:"date"`
	Text    string `db:"text"`
	PhotoID int    `db:"photo_id"`
}

type User struct {
	ID       int
	AvatarID int
}

type Photo struct {
	ID  int
	URL string
}

type Store struct {
	db *sqlx.DB
}

func (s *Store) idsStr(ids []int) string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.Itoa(id)
	}

	return strings.Join(result, ",")
}

func (s *Store) GetFeed() ([]int, error) {
	var postIds []int
	err := s.db.Select(&postIds, "select id from post order by id desc limit 30")
	if err != nil {
		return nil, fmt.Errorf("error selecting post ids")
	}

	return postIds, nil
}

func (s *Store) GetPosts(ids []int) ([]*Post, error) {
	var posts []*Post
	err := s.db.Select(&posts, "select * from post where id in ("+s.idsStr(ids)+")")
	if err != nil {
		return nil, fmt.Errorf("error selecting post ids")
	}

	return posts, err
}
