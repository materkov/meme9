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

func (o *ObjectStore) AssocAdd(id1, id2, assocType int, data *StoredAssoc) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	_, err = o.db.Exec("insert into assoc(id1, id2, type, data) values (?, ?, ?, ?)", id1, id2, assocType, dataBytes)
	if err != nil {
		return fmt.Errorf("error saving row: %w")
	}

	return nil
}

func (o *ObjectStore) AssocCount(id, assocType int) (int, error) {
	log.Printf("[INFO] AssocCount %d --(%d)--> COUNT()", id, assocType)

	cnt := 0
	err := o.db.QueryRow("select count(*) from assoc where id1 = ? and type = ?", id, assocType).Scan(&cnt)
	return cnt, err
}

func (o *ObjectStore) AssocDelete(id1, id2, assocType int) error {
	log.Printf("[INFO] AssocDelete %d --(%d)--> %d", id1, assocType, id2)

	_, err := o.db.Exec("delete from assoc where id1 = ? and id2 = ? and type = ?", id1, id2, assocType)
	if err != nil {
		return fmt.Errorf("error saving row: %w")
	}

	return nil
}

func (o *ObjectStore) AssocGet(id1, assocType, id2 int) (*StoredAssoc, error) {
	log.Printf("[INFO] AssocGet %d --(%d)--> %d", id1, assocType, id2)

	var data []byte
	err := o.db.QueryRow("select data from assoc where id1 = ? and id2 = ? and type = ?", id1, id2, assocType).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error selecting row: %w", err)
	}

	assoc := &StoredAssoc{}
	err = json.Unmarshal(data, &assoc)
	if err != nil {
	    return nil, fmt.Errorf("error unmarshaling object: %w", err)
	}

	return assoc, nil
}
