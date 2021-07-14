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

func (o *ObjectStore) ObjUpdate(object *StoredObject) error {
	objectBytes, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("error marshaling object: %w", err)
	}

	_, err = o.db.Exec("update object set data = ? where id = ?", objectBytes, object.ID)
	if err != nil {
		return fmt.Errorf("error saving to mysql: %w", err)
	}

	return nil
}

func (o *ObjectStore) AssocAdd(id1, id2 int, assocType string, data *StoredAssoc) error {
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

func (o *ObjectStore) AssocCount(id int, assocType string) (int, error) {
	log.Printf("[INFO] AssocCount %d --(%s)--> COUNT()", id, assocType)

	cnt := 0
	err := o.db.QueryRow("select count(*) from assoc where id1 = ? and type = ?", id, assocType).Scan(&cnt)
	return cnt, err
}

func (o *ObjectStore) AssocDelete(id1, id2 int, assocType string) error {
	log.Printf("[INFO] AssocDelete %d --(%s)--> %d", id1, assocType, id2)

	_, err := o.db.Exec("delete from assoc where id1 = ? and id2 = ? and type = ?", id1, id2, assocType)
	if err != nil {
		return fmt.Errorf("error saving row: %w")
	}

	return nil
}

func (o *ObjectStore) AssocGet(id1 int, assocType string, id2 int) (*StoredAssoc, error) {
	log.Printf("[INFO] AssocGet %d --(%s)--> %d", id1, assocType, id2)

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

func (o *ObjectStore) AssocRange(id1 int, assocType string, limit int) ([]*StoredAssoc, error) {
	rows, err := o.db.Query("select data from assoc where id1 = ? and type = ? order by id desc limit ?", id1, assocType, limit)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}
	defer rows.Close()

	result := make([]*StoredAssoc, 0)
	for rows.Next() {
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		obj := &StoredAssoc{}
		err = json.Unmarshal(data, obj)
		if err != nil {
			return nil, fmt.Errorf("error umarshaling assoc: %w", err)
		}

		result = append(result, obj)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return result, nil
}

func (o *ObjectStore) GenerateNextID() (int, error) {
	result, err := o.db.Exec("insert into objects(object_type) values (0)")
	if err != nil {
		return 0, fmt.Errorf("error inserting object row: %s", err)
	}

	id, _ := result.LastInsertId()
	return int(id), err
}
