package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/store"
)

type SqlUserStore struct {
	DB *sql.DB
}

func (u *SqlUserStore) Get(ids []int) (map[int]*store.User, error) {
	objectBytes, err := LoadObjects(u.DB, ids)
	if err != nil {
		return nil, err
	}

	result := map[int]*store.User{}
	for objectID, objectBytes := range objectBytes {
		object := store.User{}
		err = json.Unmarshal(objectBytes, &object)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling object: %w", err)
		}
		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (u *SqlUserStore) Add(object *store.User) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	result, err := u.DB.Exec("insert into objects(obj_type, data) values (?, ?)", store.ObjTypeUser, objectBytes)
	if err != nil {
		return err
	}

	objectID, _ := result.LastInsertId()
	object.ID = int(objectID)

	return nil
}

func (u *SqlUserStore) Update(object *store.User) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec("update objects set data = ? where id = ?", objectBytes, object.ID)
	return err
}

type UserStore interface {
	Get(ids []int) (map[int]*store.User, error)
	Add(object *store.User) error
	Update(object *store.User) error
}

type MockUserStore struct {
	nextID  int
	Objects map[int]*store.User
}

func (m *MockUserStore) Get(ids []int) (map[int]*store.User, error) {
	result := map[int]*store.User{}

	for _, objectID := range ids {
		result[objectID] = m.Objects[objectID]
	}

	return result, nil
}

func (m *MockUserStore) Add(object *store.User) error {
	m.nextID++

	object.ID = m.nextID
	m.Objects[m.nextID] = object
	return nil
}

func (m *MockUserStore) Update(object *store.User) error {
	m.Objects[object.ID] = object
	return nil
}
