package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlPollAnswerStore struct {
	DB *sql.DB
}

func (u *SqlPollAnswerStore) Get(ids []int) (map[int]*store.PollAnswer, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToString(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.PollAnswer{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.PollAnswer{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (u *SqlPollAnswerStore) Add(object *store.PollAnswer) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	_, err = u.DB.Exec("insert into objects(obj_type, data) values (?, ?)", store.ObjTypePollAnswer, objectBytes)
	return err
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
