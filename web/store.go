package main

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate go run cmd/sqlgenerator/main.go
const (
	ObjectTypePost     = 1
	ObjectTypeUser     = 2
	ObjectTypeToken    = 3
	ObjectTypePhoto    = 4
	ObjectTypeLike     = 5
	ObjectTypeComment  = 6
	ObjectTypeAPILog   = 7
	ObjectTypeFollower = 8
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

func GetIdFromToken(token string) int {
	parts := strings.Split(token, "-")
	if len(parts) < 2 {
		return 0
	}
	tokenID, _ := strconv.Atoi(parts[0])
	if tokenID <= 0 {
		return 0
	}

	return tokenID
}

type Photo struct {
	ID     int    `db:"id"`
	UserID int    `db:"user_id"`
	Path   string `db:"path"`
}

type Likes struct {
	ID     int `db:"id"`
	PostID int `db:"post_id"`
	UserID int `db:"user_id"`
	Time   int `db:"time"`
}

var store *Store

func (s *PostStore) GetByUsers(userIds []int) ([]int, error) {
	postIds, err := scanIdsList(s.db, fmt.Sprintf("select id from post where user_id in (%s) order by id desc limit 30", idsStr(userIds)))
	return postIds, err
}

func (s *UserStore) GetByVkID(vkID int) (int, error) {
	userIds, err := scanIdsList(s.db, "select id from user where vk_id = "+strconv.Itoa(vkID))
	if err != nil {
		return 0, err
	} else if len(userIds) == 0 {
		return 0, nil
	} else {
		return userIds[0], err
	}
}

func (s *UserStore) UpdateNameAvatar(user *User) error {
	_, err := s.db.Exec(
		"update user set name = ?, vk_avatar = ? where id = ?",
		user.Name, user.VkAvatar, user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating name and avatar: %s", err)
	}

	return nil
}

func (s *LikesStore) Delete(postID, userID int) error {
	_, err := s.db.Exec(
		"delete from likes where post_id = ? and user_id = ?",
		postID, userID,
	)
	return err
}

func (s *LikesStore) GetCount(postIds []int) (map[int]int, error) {
	rows, err := s.db.Query("select post_id, count(*) from likes where post_id in (" + idsStr(postIds) + ") group by post_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]int{}
	postID, count := 0, 0
	for rows.Next() {
		err := rows.Scan(&postID, &count)
		if err != nil {
			return nil, err
		}

		result[postID] = count
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *LikesStore) GetIsLiked(postIds []int, userID int) (map[int]bool, error) {
	result := map[int]bool{}

	postIds, err := scanIdsList(
		s.db,
		fmt.Sprintf("select post_id from likes where user_id = %d and post_id in (%s)", userID, idsStr(postIds)),
	)
	if err != nil {
		return nil, err
	}

	for _, postID := range postIds {
		result[postID] = true
	}

	return result, nil
}

func (s *FollowersStore) GetFollowing(userID int) ([]int, error) {
	return scanIdsList(s.db, "select user2_id from followers where user1_id = "+strconv.Itoa(userID))
}

func (s *FollowersStore) Unfollow(userFrom, userTo int) error {
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

func (s *CommentStore) GetCommentsCounts(postIds []int) (map[int]int, error) {
	rows, err := s.db.Query(fmt.Sprintf("select post_id, count(*) as cnt from comment where post_id in (%s) group by post_id", idsStr(postIds)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]int{}
	postID, count := 0, 0
	for rows.Next() {
		err = rows.Scan(&postID, &count)
		if err != nil {
			return nil, err
		}

		if count > 0 {
			result[postID] = count
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CommentStore) GetLatest(postIds []int) (map[int]int, error) {
	rows, err := s.db.Query(fmt.Sprintf("select post_id, max(id) from comment where post_id in (%s) group by post_id", idsStr(postIds)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]int{}
	postID, commentID := 0, 0
	for rows.Next() {
		err = rows.Scan(&postID, &commentID)
		if err != nil {
			return nil, err
		}

		result[postID] = commentID
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CommentStore) GetByPost(postID int) ([]int, error) {
	return scanIdsList(s.db, "select id from comment where post_id = "+strconv.Itoa(postID)+" order by id desc limit 100")
}

func (s *Store) GenerateNextID(objectType int) (int, error) {
	result, err := s.db.Exec("insert into objects(object_type) values (?)", objectType)
	if err != nil {
		return 0, fmt.Errorf("error inserting object row: %s", err)
	}

	id, _ := result.LastInsertId()
	return int(id), err
}

type APILog struct {
	ID       int    `db:"id"`
	UserID   int    `db:"user_id"`
	Method   string `db:"method"`
	Request  string `db:"request"`
	Response string `db:"response"`
}

type Followers struct {
	ID         int `db:"id"`
	User1ID    int `db:"user1_id"`
	User2ID    int `db:"user2_id"`
	FollowDate int `db:"follow_date"`
}
