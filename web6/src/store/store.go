package store

import (
	"database/sql"
	"fmt"
)

// Reserved: -2, -3, -4
const (
	FakeObjPostedPost = -1
	FakeObjConfig     = -5
)

// Reserved: 1
const (
	ObjTypeConfig     = 2
	ObjTypeUser       = 3
	ObjTypePost       = 4
	ObjTypeToken      = 5
	ObjTypePoll       = 6
	ObjTypePollAnswer = 7
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

	EdgeTypeVoted = 10
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
