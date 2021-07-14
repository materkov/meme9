package main

import (
	"fmt"
	"strconv"
	"strings"

	storeAlias "github.com/materkov/meme9/web/store"
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


type Likes struct {
	ID     int `db:"id"`
	PostID int `db:"post_id"`
	UserID int `db:"user_id"`
	Time   int `db:"time"`
}

var objectStore *storeAlias.ObjectStore
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
