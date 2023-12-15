package store2

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"strings"
)

type NodeStore interface {
	Get(ids []int) (map[int][]byte, error)
	Add(objType int, data interface{}) (int, error)
	Update(id int, data interface{}) error
}

type SqlNodes struct {
	db *sql.DB
}

func NewSqlNodes(db *sql.DB) *SqlNodes {
	return &SqlNodes{db: db}
}

func (s *SqlNodes) Update(id int, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("update objects set data = ? where id = ?", dataBytes, id)
	if err != nil {
		return fmt.Errorf("error updating objects row: %w", err)
	}

	return nil
}

func (s *SqlNodes) Get(ids []int) (map[int][]byte, error) {
	rows, err := s.db.Query(fmt.Sprintf("select id, data from objects where id in (%s)", strings.Join(utils.IdsToStrings(ids), ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[int][]byte{}
	for rows.Next() {
		objectID := 0
		var data []byte

		err := rows.Scan(&objectID, &data)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		result[objectID] = data
	}

	return result, nil
}

func (s *SqlNodes) Add(objType int, data interface{}) (int, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	result, err := s.db.Exec("insert into objects(obj_type, data) values (?, ?)", objType, dataBytes)
	if err != nil {
		return 0, fmt.Errorf("error inserting objects row: %w", err)
	}

	objectID, _ := result.LastInsertId()
	return int(objectID), nil
}

type MockNodes struct {
	nextID int
	rows   map[int][]byte
}

func NewMockNodes() *MockNodes {
	return &MockNodes{
		rows: map[int][]byte{},
	}
}

func (m *MockNodes) Get(ids []int) (map[int][]byte, error) {
	result := map[int][]byte{}
	for _, id := range ids {
		data := m.rows[id]
		if data != nil {
			result[id] = data
		}
	}

	return result, nil
}

func (m *MockNodes) Add(objType int, data interface{}) (int, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	m.nextID++
	m.rows[m.nextID] = dataBytes

	return m.nextID, nil
}

func (m *MockNodes) Update(id int, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	m.rows[id] = dataBytes
	return nil
}
