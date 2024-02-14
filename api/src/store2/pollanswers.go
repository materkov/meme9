package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/store"
)

type SqlPollAnswerStore struct {
	DB *sql.DB
}

func (u *SqlPollAnswerStore) Get(ids []int) (map[int]*store.PollAnswer, error) {
	objectBytes, err := LoadObjects(u.DB, ids)
	if err != nil {
		return nil, err
	}

	result := map[int]*store.PollAnswer{}
	for objectID, objectBytes := range objectBytes {
		object := store.PollAnswer{}
		err = json.Unmarshal(objectBytes, &object)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling object: %w", err)
		}
		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (u *SqlPollAnswerStore) Add(object *store.PollAnswer) error {
	objectID, err := AddObject(u.DB, object, store.ObjTypePollAnswer)
	if err != nil {
		return err
	}

	object.ID = objectID
	return nil
}

func (u *SqlPollAnswerStore) Update(object *store.PollAnswer) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec("update objects set data = ? where id = ?", objectBytes, object.ID)
	return err
}

type PollAnswerStore interface {
	Get(ids []int) (map[int]*store.PollAnswer, error)
	Add(object *store.PollAnswer) error
	Update(object *store.PollAnswer) error
}

type MockPollAnswerStore struct {
	nextID  int
	Objects map[int]*store.PollAnswer
}

func (m *MockPollAnswerStore) Get(ids []int) (map[int]*store.PollAnswer, error) {
	result := map[int]*store.PollAnswer{}

	for _, objectID := range ids {
		result[objectID] = m.Objects[objectID]
	}

	return result, nil
}

func (m *MockPollAnswerStore) Add(object *store.PollAnswer) error {
	m.nextID++

	object.ID = m.nextID
	m.Objects[m.nextID] = object
	return nil
}

func (m *MockPollAnswerStore) Update(object *store.PollAnswer) error {
	m.Objects[object.ID] = object
	return nil
}
