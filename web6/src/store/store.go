package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

// Reserved: -2, -3, -4
const (
	FakeObjPostedPost = -1
	FakeObjConfig     = -5
)

// Reserved: 1
const (
	ObjTypeConfig = 2
	ObjTypeUser   = 3
	ObjTypePost   = 4
	ObjTypeToken  = 5
)

var SqlClient *sql.DB

var (
	ErrObjectNotFound = fmt.Errorf("object not found")
	ErrUniqueNotFound = fmt.Errorf("unique row not found")
	ErrDuplicateEdge  = fmt.Errorf("edge duplicate")
)

func getObject(id int, objType int, obj interface{}) error {
	var data []byte
	err := SqlClient.QueryRow("select data from objects where id = ? and obj_type = ?", id, objType).Scan(&data)
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

func UpdateObject(object interface{}, id int) error {
	data, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("error marshaling to json: %w", err)
	}

	_, err = SqlClient.Exec("update objects set data = ? where id = ?", data, id)
	if err != nil {
		return fmt.Errorf("error updating row: %w", err)
	}

	return nil
}

func AddObject(objType int, object interface{}) (int, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return 0, fmt.Errorf("error marshaling to json: %w", err)
	}

	res, err := SqlClient.Exec("insert into objects(obj_type, data) values (?, ?)", objType, data)
	if err != nil {
		return 0, fmt.Errorf("error inserting row: %w", err)
	}

	objId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last id from mysql: %w", err)
	}

	return int(objId), nil
}

// Reserved: 2, 3, 4, 5
const (
	EdgeTypePosted     = 1
	EdgeTypePostedPost = 3

	EdgeTypeFollowing  = 7
	EdgeTypeFollowedBy = 8
)

type Edge struct {
	FromID int
	ToID   int
	Date   int
}

func GetEdges(fromID int, edgeType int) ([]Edge, error) {
	rows, err := SqlClient.Query("select to_id, date from edges where from_id = ? and edge_type = ? order by id desc", fromID, edgeType)
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

func GetToId(edges []Edge) []int {
	result := make([]int, len(edges))
	for i, edge := range edges {
		result[i] = edge.ToID
	}
	return result
}

var ErrNoEdge = fmt.Errorf("no edge")

func GetEdge(fromID, toID, edgeType int) (Edge, error) {
	e := Edge{
		FromID: fromID,
		ToID:   toID,
	}

	err := SqlClient.QueryRow("select date from edges where from_id = ? and to_id = ? and edge_type = ?", fromID, toID, edgeType).Scan(&e.Date)
	if errors.Is(err, sql.ErrNoRows) {
		return Edge{}, ErrNoEdge
	} else if err != nil {
		return Edge{}, err
	}

	return e, nil
}

func AddEdge(fromID, toID, edgeType int) error {
	_, err := SqlClient.Exec(
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

func DelEdge(fromID, toID, edgeType int) error {
	_, err := SqlClient.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", fromID, edgeType, toID)
	if err != nil {
		return fmt.Errorf("error deleteing edge: %w", err)
	}

	return nil
}

const (
	UniqueTypeEmail     = 1
	UniqueTypeVKID      = 2
	UniqueTypeAuthToken = 3
)

func AddUnique(keyType int, key string, objectID int) error {
	_, err := SqlClient.Exec("insert into uniques(type, `key`, object_id) values (?, ?, ?)", keyType, key, objectID)
	if err != nil {
		return fmt.Errorf("error inserting unique row: %w", err)
	}
	return nil
}

func GetUnique(keyType int, key string) (int, error) {
	if key == "" {
		return 0, ErrUniqueNotFound
	}

	objectID := 0
	err := SqlClient.QueryRow("select object_id from uniques where type = ? and `key` = ?", keyType, key).Scan(&objectID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrUniqueNotFound
	} else if err != nil {
		return 0, fmt.Errorf("error selecing unique row: %w", err)
	}

	return objectID, nil
}
