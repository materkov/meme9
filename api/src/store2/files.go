package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/api/src/store"
)

type SqlFileStore struct {
	DB *sql.DB
}

func (s *SqlFileStore) Get(ids []int) (map[int]*store.File, error) {
	objectBytes, err := LoadObjects(s.DB, ids)
	if err != nil {
		return nil, err
	}

	result := map[int]*store.File{}
	for objectID, objectBytes := range objectBytes {
		object := store.File{}
		err = json.Unmarshal(objectBytes, &object)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling object: %w", err)
		}
		object.ID = objectID
		result[objectID] = &object
	}

	return result, nil
}

func (s *SqlFileStore) Add(object *store.File) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	result, err := s.DB.Exec("insert into objects(obj_type, data) values (?, ?)", store.ObjTypeFile, objectBytes)
	if err != nil {
		return err
	}

	objectID, _ := result.LastInsertId()
	object.ID = int(objectID)

	return nil
}
