package store2

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = fmt.Errorf("unique not found")

type UniqueStore interface {
	Get(uniqType int, val string) (int, error)
	Add(uniqType int, val string, objectID int) error
}

type SqlUniqueStore struct {
	DB *sql.DB
}

func (u *SqlUniqueStore) Add(uniqType int, val string, objectID int) error {
	_, err := u.DB.Exec("insert into uniques (type, `key`, object_id) values (?, ?, ?)", uniqType, val, objectID)
	if err != nil {
		return fmt.Errorf("error inserting unique row: %w", err)
	}

	return nil
}

func (u *SqlUniqueStore) Get(uniqType int, val string) (int, error) {
	objectID := 0
	err := u.DB.QueryRow("select object_id from uniques where `type` = ? and `key` = ?", uniqType, val).Scan(&objectID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, fmt.Errorf("error selecting uniques row: %w", err)
	}

	return objectID, nil
}

type MockUniqueStore struct {
	Rows map[string]int
}

func (m *MockUniqueStore) Add(uniqType int, val string, objectID int) error {
	key := fmt.Sprintf("%d:%s", uniqType, val)

	_, ok := m.Rows[key]
	if ok {
		return fmt.Errorf("duplicate id: %d-%s", objectID, val)
	}

	m.Rows[key] = objectID
	return nil
}

func (m *MockUniqueStore) Get(uniqType int, val string) (int, error) {
	key := fmt.Sprintf("%d:%s", uniqType, val)
	objectID, ok := m.Rows[key]
	if !ok {
		return 0, ErrNotFound
	}

	return objectID, nil
}
