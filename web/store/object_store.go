package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type ObjectStore struct {
	db *sql.DB
}

func NewObjectStore(db *sql.DB) *ObjectStore {
	return &ObjectStore{db: db}
}

func (o *ObjectStore) ObjGet(id int) (*StoredObject, error) {
	log.Printf("[INFO] ObjGet: id %d", id)

	var data []byte
	err := o.db.QueryRow("select data from object where id = " + strconv.Itoa(id)).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error selecting row: %s", err)
	}

	obj := &StoredObject{}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling object: %w", err)
	}

	return obj, nil
}

func (o *ObjectStore) ObjAdd(object *StoredObject) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("error marshaling object: %w", err)
	}

	_, err = o.db.Exec("insert into object(id, type, data) values (?, ?, ?)", object.ID, 0, objectBytes)
	if err != nil {
		return fmt.Errorf("error saving to mysql: %w", err)
	}

	return nil
}
