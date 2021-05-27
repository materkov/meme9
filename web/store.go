package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	ObjectTypePost  = 1
	ObjectTypeUser  = 2
	ObjectTypeToken = 3
)

type Post struct {
	ID      int    `db:"id"`
	UserID  int    `db:"user_id"`
	Date    int    `db:"date"`
	Text    string `db:"text"`
	PhotoID int    `db:"photo_id"`
}

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	AvatarID int    `db:"avatar_id"`
	VkID     int    `db:"vk_id"`
	VkAvatar string `db:"vk_avatar"`
}

type Token struct {
	ID     int    `db:"id"`
	Token  string `db:"token"`
	UserID int    `db:"user_id"`
}

type Store struct {
	db *sqlx.DB
}

var store Store

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

func (s *Store) GetPostsByUsers(userIds []int) ([]int, error) {
	var postIds []int
	err := s.db.Select(
		&postIds,
		fmt.Sprintf("select id from post where user_id in (%s) order by id desc limit 30", s.idsStr(userIds)),
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting post ids")
	}

	return postIds, nil
}

func (s *Store) GetPosts(ids []int) ([]*Post, error) {
	var posts []*Post
	err := s.db.Select(&posts, "select id, user_id, date, text, COALESCE(photo_id, 0) as photo_id from post where id in ("+s.idsStr(ids)+")")
	if err != nil {
		return nil, fmt.Errorf("error selecting post ids")
	}

	return posts, err
}

func (s *Store) AddPost(post *Post) error {
	_, err := s.db.Exec(
		"insert into post (id, user_id, date, text) values (?, ?, ?, ?)",
		post.ID, post.UserID, post.Date, post.Text,
	)
	if err != nil {
		return fmt.Errorf("error inserting post row: %w", err)
	}

	return nil
}

func (s *Store) GetByVkID(vkID int) (int, error) {
	userID := 0
	err := s.db.Get(&userID, "select id from user where vk_id = ?", vkID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("error selecting user by vk id: %w", err)
	}

	return userID, err
}

func (s *Store) UpdateNameAvatar(user *User) error {
	_, err := s.db.Exec(
		"update user set name = ?, vk_avatar = ? where id = ?",
		user.Name, user.VkAvatar, user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating name and avatar: %s", err)
	}

	return nil
}

func (s *Store) GetToken(tokenStr string) (*Token, error) {
	token := Token{}
	err := s.db.Get(&token, "select * from token where token = ?", tokenStr)
	if err != nil {
		return nil, fmt.Errorf("error selecting token row: %w", err)
	}

	return &token, err
}

func (s *Store) AddToken(token *Token) error {
	result, err := s.db.Exec(
		"insert into token(token, user_id) values (?, ?)",
		token.Token, token.UserID,
	)
	if err != nil {
		return fmt.Errorf("error inserting token row: %w", err)
	}

	id, _ := result.LastInsertId()
	token.ID = int(id)

	return nil
}

func (s *Store) GetUsers(ids []int) ([]*User, error) {
	var users []*User
	err := s.db.Select(&users, "select id, coalesce(name, '') as name, coalesce(vk_avatar, '') as vk_avatar, coalesce(vk_id, 0) as vk_id from user where id in ("+s.idsStr(ids)+")")
	if err != nil {
		return nil, fmt.Errorf("error selecting users: %w", err)
	}

	return users, err
}

func (s *Store) AddUserByVK(user *User) error {
	_, err := s.db.Exec("insert into user(id, vk_id) values (?, ?)", user.ID, user.VkID)
	return err
}

func (s *Store) AddLike(postID, userID int) error {
	_, err := s.db.Exec(
		"insert into likes(post_id, user_id, time) values (?, ?, ?)",
		postID, userID, time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("error inserting like: %s", err)
	}

	return nil
}

func (s *Store) DeleteLike(postID, userID int) error {
	_, err := s.db.Exec(
		"delete from likes where post_id = ? and user_id = ?",
		postID, userID,
	)
	if err != nil {
		return fmt.Errorf("error inserting like: %s", err)
	}

	return nil
}

func (s *Store) GetLikesCount(postIds []int) (map[int]int, error) {
	result := map[int]int{}

	var rows []struct {
		PostID int `db:"post_id"`
		Count  int `db:"count"`
	}

	err := s.db.Select(
		&rows,
		"select post_id as post_id, count(*) as count from likes where post_id in ("+s.idsStr(postIds)+") group by post_id",
	)
	if err != nil {
		return result, fmt.Errorf("error selecting likes count: %s", err)
	}

	for _, row := range rows {
		result[row.PostID] = row.Count
	}

	return result, nil
}

func (s *Store) GetIsLiked(postIds []int, userID int) (map[int]bool, error) {
	result := map[int]bool{}

	var likedPosts []int
	err := s.db.Select(
		&likedPosts,
		fmt.Sprintf("select post_id from likes where user_id = %d and post_id in (%s)", userID, s.idsStr(postIds)),
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting likes count: %s", err)
	}

	for _, postID := range likedPosts {
		result[postID] = true
	}

	return result, nil
}

func (s *Store) GetFollowing(userID int) ([]int, error) {
	var userIds []int
	err := s.db.Select(&userIds, "select user2_id from followers where user1_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error selecting followers: %w", err)
	}

	return userIds, err
}

func (s *Store) Follow(userFrom, userTo int) error {
	_, err := s.db.Exec(
		"insert into followers (user1_id, user2_id, follow_date) values (?, ?, ?)",
		userFrom, userTo, time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("error inserting followers row: %w", err)
	}

	return nil
}

func (s *Store) Unfollow(userFrom, userTo int) error {
	_, err := s.db.Exec(
		"delete from followers where user1_id = ? and user2_id = ?",
		userFrom, userTo,
	)
	if err != nil {
		return fmt.Errorf("error deleting followers row: %w", err)
	}

	return nil
}

type Comment struct {
	ID     int    `db:"id"`
	PostID int    `db:"post_id"`
	UserID int    `db:"user_id"`
	Text   string `db:"text"`
	Date   int    `db:"date"`
}

func (s *Store) AddComment(comment *Comment) error {
	result, err := s.db.Exec(
		"insert into comment(post_id, text, date, user_id) values (?, ?, ?, ?)",
		comment.PostID, comment.Text, comment.Date, comment.UserID,
	)
	if err != nil {
		return fmt.Errorf("error inserting row: %w", err)
	}

	id, _ := result.LastInsertId()
	comment.ID = int(id)

	return nil
}

func (s *Store) GetCommentsCounts(postIds []int) (map[int]int, error) {
	var rows []struct {
		PostID int `db:"post_id"`
		Count  int `db:"cnt"`
	}
	err := s.db.Select(
		&rows,
		fmt.Sprintf("select post_id, count(*) as cnt from comment where post_id in (%s) group by post_id", s.idsStr(postIds)),
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting comment rows: %w", err)
	}

	result := map[int]int{}
	for _, row := range rows {
		if row.Count > 0 {
			result[row.PostID] = row.Count
		}
	}

	return result, nil
}

func (s *Store) GetLatestComments(postIds []int) (map[int]int, error) {
	var rows []struct {
		PostID    int `db:"post_id"`
		CommentID int `db:"comment_id"`
	}
	err := s.db.Select(
		&rows,
		fmt.Sprintf("select post_id, max(id) as comment_id from comment where post_id in (%s) group by post_id", s.idsStr(postIds)),
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting comment rows: %w", err)
	}

	result := map[int]int{}
	for _, row := range rows {
		result[row.PostID] = row.CommentID
	}

	return result, nil
}

func (s *Store) GetCommentsByPost(postID int) ([]int, error) {
	var commentIds []int
	err := s.db.Select(
		&commentIds,
		"select id from comment where post_id = ? order by id desc limit 100",
		postID,
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting comment rows: %w", err)
	}

	return commentIds, err
}

func (s *Store) GetComments(ids []int) ([]*Comment, error) {
	var comments []*Comment
	err := s.db.Select(
		&comments,
		fmt.Sprintf("select * from comment where id in (%s)", s.idsStr(ids)),
	)
	if err != nil {
		return nil, fmt.Errorf("error selecting comment rows: %w", err)
	}

	return comments, err
}

func (s *Store) GenerateNextID(objectType int) (int, error) {
	result, err := s.db.Exec("insert into objects(object_type) values (?)", objectType)
	if err != nil {
		return 0, fmt.Errorf("error inserting object row: %s", err)
	}

	id, _ := result.LastInsertId()
	return int(id), err
}
