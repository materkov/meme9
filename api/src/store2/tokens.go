package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/store"
)

type SqlTokenStore struct {
	DB *sql.DB
}

func (u *SqlTokenStore) Get(ids []int) (map[int]*store.Token, error) {
	objectBytes, err := LoadObjects(u.DB, ids)
	if err != nil {
		return nil, err
	}

	result := map[int]*store.Token{}
	for objectID, objectBytes := range objectBytes {
		object := store.Token{}
		err = json.Unmarshal(objectBytes, &object)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling object: %w", err)
		}
		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (u *SqlTokenStore) Add(object *store.Token) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	result, err := u.DB.Exec("insert into objects(obj_type, data) values (?, ?)", store.ObjTypeToken, objectBytes)
	if err != nil {
		return err
	}

	objectID, _ := result.LastInsertId()
	object.ID = int(objectID)

	return nil
}

type TokenStore interface {
	Get(ids []int) (map[int]*store.Token, error)
	Add(object *store.Token) error
}

type MockTokenStore struct {
	nextID  int
	Objects map[int]*store.Token
}

func (m *MockTokenStore) Get(ids []int) (map[int]*store.Token, error) {
	result := map[int]*store.Token{}

	for _, objectID := range ids {
		result[objectID] = m.Objects[objectID]
	}

	return result, nil
}

func (m *MockTokenStore) Add(object *store.Token) error {
	m.nextID++

	object.ID = m.nextID
	m.Objects[m.nextID] = object
	return nil
}
