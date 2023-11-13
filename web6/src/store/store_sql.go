package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

type Store interface {
	getObject(id int, objType int, obj interface{}) error
	AddObject(objType int, object interface{}) (int, error)
	UpdateObject(object interface{}, id int) error

	AddEdge(fromID, toID, edgeType int) error
	GetEdge(fromID, toID, edgeType int) (*Edge, error)
	CountEdges(fromID, edgeType int) (int, error)
	GetEdges(fromID int, edgeType int) ([]Edge, error)
	DelEdge(fromID, toID, edgeType int) error
}

type SqlStore struct {
	DB *sql.DB
}

func (s *SqlStore) getObject(id int, objType int, obj interface{}) error {
	var data []byte
	err := s.DB.QueryRow("select data from objects where id = ? and obj_type = ?", id, objType).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrObjectNotFound
	} else if err != nil {
		return fmt.Errorf("error selecting database: %w", err)
	}

	err = json.Unmarshal(data, obj)
	if err != nil {
		return fmt.Errorf("error unmarshaling object %d: %w", objType, err)
	}

	return nil

}

func (s *SqlStore) AddObject(objType int, object interface{}) (int, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return 0, fmt.Errorf("error marshaling to json: %w", err)
	}

	res, err := s.DB.Exec("insert into objects(obj_type, data) values (?, ?)", objType, data)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %w", err)
	}

	objId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last id from mysql: %w", err)
	}

	return int(objId), nil
}

func (s *SqlStore) AddEdge(fromID, toID, edgeType int) error {
	_, err := s.DB.Exec(
		"insert into edges(from_id, to_id, edge_type, date) values (?, ?, ?, ?)",
		fromID, toID, edgeType, time.Now().Unix(),
	)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
			return ErrDuplicateEdge
		}
		return fmt.Errorf("error inserting edge: %s", err)
	}

	return nil
}

func (s *SqlStore) GetEdge(fromID, toID, edgeType int) (*Edge, error) {
	date := 0
	err := s.DB.QueryRow("select date from edges where from_id = ? and to_id = ? and edge_type = ?", fromID, toID, edgeType).Scan(&date)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEdge
	} else if err != nil {
		return nil, err
	}

	e := &Edge{
		FromID: fromID,
		ToID:   toID,
		Date:   date,
	}
	return e, nil
}

func (s *SqlStore) CountEdges(fromID, edgeType int) (int, error) {
	cnt := 0
	err := s.DB.QueryRow("select count(*) from edges where from_id = ? and edge_type = ?", fromID, edgeType).Scan(&cnt)
	return cnt, err
}

func (s *SqlStore) GetEdges(fromID int, edgeType int) ([]Edge, error) {
	rows, err := s.DB.Query("select to_id, date from edges where from_id = ? and edge_type = ? order by id desc", fromID, edgeType)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}
	defer rows.Close()

	var results []Edge
	for rows.Next() {
		e := Edge{FromID: fromID}
		err = rows.Scan(&e.ToID, &e.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning edge row: %w", err)
		}

		results = append(results, e)
	}

	return results, err
}

func (s *SqlStore) DelEdge(fromID, toID, edgeType int) error {
	_, err := s.DB.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", fromID, edgeType, toID)
	if err != nil {
		return fmt.Errorf("error deleteing edge: %w", err)
	}

	return nil
}

func (s *SqlStore) UpdateObject(object interface{}, id int) error {
	data, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("error marshaling to json: %w", err)
	}

	_, err = s.DB.Exec("update objects set data = ? where id = ?", data, id)
	if err != nil {
		return fmt.Errorf("error updating row: %w", err)
	}

	return nil
}
