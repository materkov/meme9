package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
)

type SqlPollStore struct {
	DB *sql.DB
}

func (u *SqlPollStore) Get(ids []int) (map[int]*store.Poll, error) {
	rows, err := u.DB.Query(fmt.Sprintf("select id, data from objects where id in (%s)", utils.IdsToCommaSeparated(ids)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int]*store.Poll{}
	for rows.Next() {
		objectID := 0
		var data []byte
		err = rows.Scan(&objectID, &data)
		if err != nil {
			return nil, err
		}

		object := store.Poll{}
		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, err
		}

		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (u *SqlPollStore) Add(object *store.Poll) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	result, err := u.DB.Exec("insert into objects(obj_type, data) values (?, ?)", store.ObjTypePoll, objectBytes)
	if err != nil {
		return err
	}

	objectID, _ := result.LastInsertId()
	object.ID = int(objectID)

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
