package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/store"
)

type SqlPollStore struct {
	DB *sql.DB
}

func (u *SqlPollStore) Get(ids []int) (map[int]*store.Poll, error) {
	objectBytes, err := LoadObjects(u.DB, ids)
	if err != nil {
		return nil, err
	}

	result := map[int]*store.Poll{}
	for objectID, objectBytes := range objectBytes {
		object := store.Poll{}
		err = json.Unmarshal(objectBytes, &object)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling object: %w", err)
		}
		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (u *SqlPollStore) Add(object *store.Poll) error {
	objectID, err := AddObject(u.DB, object, store.ObjTypePoll)
	if err != nil {
		return err
	}

	object.ID = objectID
	return nil
}

func (u *SqlPollStore) Update(object *store.Poll) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec("update objects set data = ? where id = ?", objectBytes, object.ID)
	return err
}

type PollStore interface {
	Get(ids []int) (map[int]*store.Poll, error)
	Add(object *store.Poll) error
	Update(object *store.Poll) error
}

type MockPollStore struct {
	nextID  int
	Objects map[int]*store.Poll
}

func (m *MockPollStore) Get(ids []int) (map[int]*store.Poll, error) {
	result := map[int]*store.Poll{}

	for _, objectID := range ids {
		result[objectID] = m.Objects[objectID]
	}

	return result, nil
}

func (m *MockPollStore) Add(object *store.Poll) error {
	m.nextID++

	object.ID = m.nextID
	m.Objects[m.nextID] = object
	return nil
}

func (m *MockPollStore) Update(object *store.Poll) error {
	m.Objects[object.ID] = object
	return nil
}
