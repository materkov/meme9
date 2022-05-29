package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type ObjectType int

const (
	ObjectTypePost = 1
	ObjectTypeUser = 2
)

type SqlObjectStore struct {
	db *sql.DB
}

func (s *SqlObjectStore) GetByIdMany(ids []int) ([]interface{}, error) {
	rows, err := s.db.Query(fmt.Sprintf("select id, object_type, object from object where id in (%s)", idsStr(ids)))
	if err != nil {
		return nil, fmt.Errorf("error selectng objects: %w", err)
	}
	defer rows.Close()

	var objects []interface{}

	for rows.Next() {
		var objId, objType int
		var data []byte
		err = rows.Scan(&objId, &objType, &data)
		if err != nil {
			return nil, fmt.Errorf("error scanning object row: %w", err)
		}

		var object interface{}
		switch objType {
		case ObjectTypeUser:
			object = User{ID: objId}
		case ObjectTypePost:
			object = Post{ID: objId}
		}

		err = json.Unmarshal(data, &object)
		if err != nil {
			return nil, fmt.Errorf("error unamrshaling object: %w", err)
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func (s *SqlObjectStore) Add(objType int, obj interface{}) error {
	objSerialized, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("error marshaling object: %s", err)
	}

	objID := 0
	err = s.db.QueryRow(
		"insert into object(object_type, object) values (?, ?) returning id",
		objType, objSerialized,
	).Scan(&objID)
	if err != nil {
		return fmt.Errorf("error inserting object row: %w", err)
	}

	switch obj := obj.(type) {
	case *User:
		obj.ID = objID
	case *Post:
		obj.ID = objID
	}

	return nil
}
