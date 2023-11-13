package store

import (
	"database/sql"
	"errors"
	"fmt"
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

// Reserved: 2, 3, 4, 5
const (
	EdgeTypePosted     = 1
	EdgeTypePostedPost = 3

	EdgeTypeFollowing  = 7
	EdgeTypeFollowedBy = 8

	EdgeTypeLiked = 9
)

type Edge struct {
	FromID int
	ToID   int
	Date   int
}

func GetToId(edges []Edge) []int {
	result := make([]int, len(edges))
	for i, edge := range edges {
		result[i] = edge.ToID
	}
	return result
}

var ErrNoEdge = fmt.Errorf("no edge")

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
