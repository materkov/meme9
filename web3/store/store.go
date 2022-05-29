package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Object interface {
	ObjectID() int
}

type User struct {
	ID   int
	Name string
}

func (u *User) ObjectID() int { return u.ID }

type Post struct {
	ID     int
	UserID int
	Text   string
	Date   int
}

func (p *Post) ObjectID() int { return p.ID }

const (
	ListSubscribedTo = 1
	ListPosted       = 2

	ObjectUser = 1
	ObjectPost = 2
)

type Store struct {
	DB *sql.DB
}

func parseObject(objID int, objType int, data []byte) (Object, error) {
	if objType == ObjectUser {
		obj := User{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, fmt.Errorf("error unmarshaling user: %w", err)
		}
		obj.ID = objID
		return &obj, nil
	} else if objType == ObjectPost {
		obj := Post{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, fmt.Errorf("error unmarshaling post: %w", err)
		}
		obj.ID = objID
		return &obj, nil
	} else {
		return nil, fmt.Errorf("unmknown type: %d", objType)
	}
}

func idsList(ids []int) string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.Itoa(id)
	}

	return strings.Join(result, ",")
}

func (s *Store) ObjGet(ids []int) (map[int]Object, error) {
	rows, err := s.DB.Query("select id, data, type from object where id in (" + idsList(ids) + ")")
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}
	defer rows.Close()

	result := map[int]Object{}
	for rows.Next() {
		var data []byte
		objID := 0
		objType := 0

		err = rows.Scan(&objID, &data, &objType)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		obj, _ := parseObject(objID, objType, data)
		result[objID] = obj
	}

	return result, nil
}

func (s *Store) ListGet(objectID int, listType int) ([]int, error) {
	rows, err := s.DB.Query("select object2 from list where object1 = ? and type = ?", objectID, listType)
	if err != nil {
		return nil, fmt.Errorf("error selecting list: %w", err)
	}
	defer rows.Close()

	var result []int
	for rows.Next() {
		objectID := 0
		err = rows.Scan(&objectID)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		result = append(result, objectID)
	}

	return result, err
}

func (s *Store) ObjAdd(objectID int, objectType int, obj interface{}) error {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error marshaling object: %w", err)
	}

	_, err = s.DB.Exec("insert into object(id, type, data) values (?, ?, ?)", objectID, objectType, objBytes)
	if err != nil {
		return fmt.Errorf("error inserting db row: %w", err)
	}

	return nil
}

func (s *Store) ListAdd(object1, listType, object2 int) error {
	date := time.Now().Unix()
	_, err := s.DB.Exec("insert into list(object1, type, object2, date) values (?, ?, ?, ?)", object1, listType, object2, date)
	if err != nil {
		return fmt.Errorf("error isnerting db row: %w", err)
	}

	return nil
}

func GenerateID() int {
	return int(time.Now().Unix())
}
