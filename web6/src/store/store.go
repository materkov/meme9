package store

import (
	"database/sql"
	"fmt"
)

// Reserved: -2, -3, -4, -1

const (
	FakeObjConfig = -5
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
