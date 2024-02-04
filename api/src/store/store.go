package store

import (
	"database/sql"
)

// Reserved: 1, 2
const (
	ObjTypeUser       = 3
	ObjTypePost       = 4
	ObjTypeToken      = 5
	ObjTypePoll       = 6
	ObjTypePollAnswer = 7
)

var SqlClient *sql.DB

// Reserved: 2, 3, 4, 5
const (
	EdgeTypePosted     = 1
	EdgeTypePostedPost = 3

	EdgeTypeFollowing  = 7
	EdgeTypeFollowedBy = 8

	EdgeTypeLiked = 9

	EdgeTypeVoted = 10

	EdgeTypeBookmarked = 11
)
