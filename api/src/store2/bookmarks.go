package store2

import (
	"database/sql"
	"fmt"
	"github.com/materkov/meme9/api/src/pkg/utils"
	"github.com/materkov/meme9/api/src/store"
	"slices"
	"time"
)

type BookmarkStore interface {
	IsBookmarked(postIds []int, viewerID int) (map[int]bool, error)
	Add(postID, viewerID int) error
	Remove(postID, userID int) error
	List(userId int, after int, count int) ([]BookmarkItem, error)
}

type SqlBookmarks struct {
	DB *sql.DB
}

func (s *SqlBookmarks) IsBookmarked(postIds []int, viewerID int) (map[int]bool, error) {
	query := `
select to_id
from edges
where from_id = %d and edge_type=%d and to_id in (%s)
`
	rows, err := s.DB.Query(fmt.Sprintf(query,
		viewerID, store.EdgeTypeBookmarked, utils.IdsToCommaSeparated(postIds),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	isBookmarked := map[int]bool{}

	for rows.Next() {
		postID := 0
		err = rows.Scan(&postID)
		if err != nil {
			return nil, err
		}

		isBookmarked[postID] = true
	}

	return isBookmarked, nil
}

func (s *SqlBookmarks) Add(postID, viewerID int) error {
	_, err := s.DB.Exec("insert ignore into edges (from_id, to_id, edge_type, date) values (?, ?, ?, unix_timestamp())", viewerID, postID, store.EdgeTypeBookmarked)
	return err
}

func (l *SqlBookmarks) Remove(postID, userID int) error {
	_, err := l.DB.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", userID, store.EdgeTypeBookmarked, postID)
	return err
}

type BookmarkItem struct {
	PostID int
	Date   int
}

func (s *SqlBookmarks) List(userId int, after int, count int) ([]BookmarkItem, error) {
	query := "select to_id, date from edges where from_id = ? and edge_type = ?"
	args := []any{userId, store.EdgeTypeBookmarked}

	if after != 0 {
		query += " and date < ?"
		args = append(args, after)
	}

	query += " order by date desc limit ?"
	args = append(args, count)

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []BookmarkItem
	for rows.Next() {
		item := BookmarkItem{}
		err = rows.Scan(&item.PostID, &item.Date)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, err
}

type MockBookmarks struct {
	rows [][]int // postId, date, userId
}

func (m *MockBookmarks) IsBookmarked(postIds []int, viewerID int) (map[int]bool, error) {
	result := map[int]bool{}
	for _, postID := range postIds {
		for _, row := range m.rows {
			if row[0] == postID && row[2] == viewerID {
				result[postID] = true
			}
		}
	}

	return result, nil
}

func (m *MockBookmarks) Add(postID, viewerID int) error {
	for _, row := range m.rows {
		if row[0] == postID && row[2] == viewerID {
			return nil
		}
	}

	m.rows = append(m.rows, []int{postID, int(time.Now().Unix()), viewerID})
	return nil
}

func (m *MockBookmarks) Remove(postID, userID int) error {
	var filtered [][]int
	for _, row := range m.rows {
		if row[0] == postID && row[2] == userID {
			continue
		}
		filtered = append(filtered, row)
	}

	m.rows = filtered
	return nil
}

func (m *MockBookmarks) List(userId int, after int, count int) ([]BookmarkItem, error) {
	var result []BookmarkItem
	for _, row := range m.rows {
		if row[2] == userId {
			if after != 0 && after >= row[1] {
				continue
			}
			result = append(result, BookmarkItem{
				PostID: row[0],
				Date:   row[1],
			})
		}
	}

	slices.Reverse(result)
	if len(result) > count {
		result = result[:count]
	}

	return result, nil
}
