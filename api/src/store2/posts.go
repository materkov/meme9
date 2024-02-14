package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/store"
)

type SqlPostStore struct {
	DB *sql.DB
}

func (u *SqlPostStore) Get(ids []int) (map[int]*store.Post, error) {
	objectBytes, err := LoadObjects(u.DB, ids)
	if err != nil {
		return nil, err
	}

	result := map[int]*store.Post{}
	for objectID, objectBytes := range objectBytes {
		object := store.Post{}
		err = json.Unmarshal(objectBytes, &object)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling object: %w", err)
		}
		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (u *SqlPostStore) Add(object *store.Post) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	result, err := u.DB.Exec("insert into objects(obj_type, data) values (?, ?)", store.ObjTypePost, objectBytes)
	if err != nil {
		return err
	}

	objectID, _ := result.LastInsertId()
	object.ID = int(objectID)

	return nil
}

func (u *SqlPostStore) Update(object *store.Post) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec("update objects set data = ? where id = ?", objectBytes, object.ID)
	return err
}

type PostStore interface {
	Get(ids []int) (map[int]*store.Post, error)
	Add(object *store.Post) error
	Update(object *store.Post) error
}

type MockPostStore struct {
	nextID  int
	Objects map[int]*store.Post
}

func (m *MockPostStore) Get(ids []int) (map[int]*store.Post, error) {
	result := map[int]*store.Post{}

	for _, objectID := range ids {
		result[objectID] = m.Objects[objectID]
	}

	return result, nil
}

func (m *MockPostStore) Add(object *store.Post) error {
	m.nextID++

	object.ID = m.nextID
	m.Objects[m.nextID] = object
	return nil
}

func (m *MockPostStore) Update(object *store.Post) error {
	m.Objects[object.ID] = object
	return nil
}
