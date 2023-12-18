package store2

import (
	"database/sql"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
	"slices"
	"sort"
)

type SqlWall struct {
	DB *sql.DB
}

func (s *SqlWall) Get(userIds []int, after int, limit int) ([]int, error) {
	queryAfter := ""
	if after != 0 {
		queryAfter = fmt.Sprintf(" and to_id < %d", after)
	}

	query := fmt.Sprintf(
		"select to_id from edges where from_id in (%s) and edge_type = %d %s order by id desc limit %d",
		utils.IdsToCommaSeparated(userIds),
		store.EdgeTypePosted,
		queryAfter,
		limit,
	)

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []int
	for rows.Next() {
		postID := 0
		err = rows.Scan(&postID)
		if err != nil {
			return nil, err
		}

		result = append(result, postID)
	}

	return result, nil
}

func (s *SqlWall) Add(userID, postID int) error {
	_, err := s.DB.Exec("insert into edges (from_id, to_id, edge_type, date) values (?, ?, ?, unix_timestamp())", userID, postID, store.EdgeTypePosted)
	return err
}

func (s *SqlWall) Delete(userID, postID int) error {
	_, err := s.DB.Exec("delete from edges where from_id = ? and to_id = ? and edge_type = ?", userID, postID, store.EdgeTypePosted)
	return err
}

func (s *SqlWall) GetLatest() ([]int, error) {
	rows, err := s.DB.Query("select to_id from edges where edge_type = ? order by id desc limit 1000", store.EdgeTypePosted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIds []int
	for rows.Next() {
		postID := 0
		err = rows.Scan(&postID)
		if err != nil {
			return nil, err
		}

		postIds = append(postIds, postID)
	}

	return postIds, nil
}

type Wall interface {
	Get(userIds []int, after int, limit int) ([]int, error)
	Add(userID, postID int) error
	Delete(userID, postID int) error
	GetLatest() ([]int, error)
}

type MockWall struct {
	Posts map[int][]int
}

func (m *MockWall) Get(userIds []int, after int, limit int) ([]int, error) {
	var result []int
	for _, userID := range userIds {
		result = append(result, m.Posts[userID]...)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(result)))

	if after != 0 {
		idx := slices.Index(result, after)
		if idx != -1 {
			result = result[idx+1:]
		}
	}
	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func (m *MockWall) Add(userID, postID int) error {
	m.Posts[userID] = append(m.Posts[userID], postID)
	return nil
}

func (m *MockWall) Delete(userID, postID int) error {
	var newList []int
	for _, userPostID := range m.Posts[userID] {
		if userPostID != postID {
			newList = append(newList, userPostID)
		}
	}

	m.Posts[userID] = newList
	return nil
}

func (m *MockWall) GetLatest() ([]int, error) {
	var posts []int
	for _, userPosts := range m.Posts {
		posts = append(posts, userPosts...)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(posts)))
	return posts, nil
}
