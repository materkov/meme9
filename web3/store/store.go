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
	VkID int
}

type Config struct {
	ID int

	VKAppID     int
	VKAppSecret string

	AuthTokenSecret string
}

func (c *Config) ObjectID() int { return c.ID }

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

	ObjectUser   = 1
	ObjectPost   = 2
	ObjectConfig = 3

	ObjectIDConfig = 1
)

type SqlStore struct {
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
	} else if objType == ObjectConfig {
		obj := Config{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return nil, fmt.Errorf("error unmarshaling config: %w", err)
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

func (s *SqlStore) ObjGet(ids []int) (map[int]Object, error) {
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

func (s *SqlStore) ListGet(objectID int, listType int) ([]int, error) {
	rows, err := s.DB.Query("select object2 from list where object1 = ? and type = ? order by id desc", objectID, listType)
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

func (s *SqlStore) ObjAdd(objectID int, objectType int, obj interface{}) error {
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

func (s *SqlStore) ListAdd(object1, listType, object2 int) error {
	date := time.Now().Unix()
	_, err := s.DB.Exec("insert into list(object1, type, object2, date) values (?, ?, ?, ?)", object1, listType, object2, date)
	if err != nil {
		return fmt.Errorf("error isnerting db row: %w", err)
	}

	return nil
}

func (s *SqlStore) ListCount(objectID, listType int) (int, error) {
	count := 0
	err := s.DB.QueryRow("select count(*) from list where object1 = ? and type = ?", objectID, listType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("sql error: %w", err)
	}

	return count, nil
}

func GenerateID() int {
	return int(time.Now().Unix())
}

const (
	MappingVKID = 1
)

func (s *SqlStore) GetMapping(keyType int, key string) (int, error) {
	objectID := 0
	err := s.DB.QueryRow("select object from mapping where `key_type` = ? and `key` = ?", keyType, key).Scan(&objectID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("error selecting mapping: %w", err)
	}

	return objectID, nil
}

func (s *SqlStore) SaveMapping(keyType int, key string, objectID int) error {
	_, err := s.DB.Exec("insert into mapping(key_type, key, object) values (?, ?, ?)", keyType, key, objectID)
	if err != nil {
		return fmt.Errorf("error saving mapping row: %w", err)
	}

	return nil
}

type Store interface {
	ListGet(objectID int, listType int) ([]int, error)
	ListAdd(object1, listType, object2 int) error
	ListCount(objectID, listType int) (int, error)

	ObjGet(ids []int) (map[int]Object, error)
	ObjAdd(objectID int, objectType int, obj interface{}) error

	GetMapping(keyType int, key string) (int, error)
	SaveMapping(keyType int, key string, objectID int) error
}
