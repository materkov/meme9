package store2

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
	"strconv"
	"strings"
)

type SqlLikes struct {
	DB *sql.DB
}

func (l *SqlLikes) Add(objectID, userID int) error {
	_, err := l.DB.Exec("insert into edges (from_id, to_id, edge_type, date) values (?, ?, ?, now())", objectID, userID, store.EdgeTypeLiked)
	return err
}

func (l *SqlLikes) Remove(objectID, userID int) error {
	_, err := l.DB.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", objectID, store.EdgeTypeLiked, userID)
	return err
}

func (l *SqlLikes) Get(ctx context.Context, postIds []int, viewerID int) (counters map[int]int, isLiked map[int]bool, err error) {
	query := `
select from_id, count(*), sum(to_id = %d)
from edges
where from_id in (%s) and edge_type=%d
group by from_id
`
	rows, err := l.DB.Query(fmt.Sprintf(query, viewerID, strings.Join(utils.IdsToStrings(postIds), ","), store.EdgeTypeLiked))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	isLiked = map[int]bool{}
	counters = map[int]int{}

	for rows.Next() {
		postID, count, isLikedInt := 0, 0, 0
		err = rows.Scan(&postID, &count, &isLikedInt)
		if err != nil {
			return nil, nil, err
		}

		counters[postID] = count
		if isLikedInt > 0 {
			isLiked[postID] = true
		}
	}

	return counters, isLiked, nil
}

type Likes interface {
	Add(objectID, userID int) error
	Remove(objectID, userID int) error
	Get(ctx context.Context, postIds []int, viewerID int) (counters map[int]int, isLiked map[int]bool, err error)
}

type MockLikes struct {
	Rows map[string]bool
}

func (m *MockLikes) Add(objectID, userID int) error {
	key := fmt.Sprintf("%d:%d", objectID, userID)
	m.Rows[key] = true
	return nil
}

func (m *MockLikes) Remove(objectID, userID int) error {
	key := fmt.Sprintf("%d:%d", objectID, userID)
	delete(m.Rows, key)
	return nil
}

func (m *MockLikes) Get(ctx context.Context, postIds []int, viewerID int) (counters map[int]int, isLiked map[int]bool, err error) {
	counters = map[int]int{}
	isLiked = map[int]bool{}

	for _, postID := range postIds {
		key := fmt.Sprintf("%d:%d", postID, viewerID)
		if m.Rows[key] {
			isLiked[postID] = true
		}
	}

	for key := range m.Rows {
		parts := strings.Split(key, ":")
		postID, _ := strconv.Atoi(parts[0])
		counters[postID]++
	}

	return counters, isLiked, err
}
