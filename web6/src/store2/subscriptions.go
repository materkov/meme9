package store2

import (
	"database/sql"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
	"strings"
)

type SqlSubscriptions struct {
	DB *sql.DB
}

func (s *SqlSubscriptions) Follow(userID, targetID int) error {
	_, err := s.DB.Exec("insert into edges(from_id, to_id, edge_type, date) values (?, ?, ?, unix_timestamp())", userID, targetID, store.EdgeTypeFollowing)
	return err
}

func (s *SqlSubscriptions) Unfollow(userID, targetID int) error {
	_, err := s.DB.Exec("delete from edges where from_id = ? and edge_type = ? and to_id = ?", userID, store.EdgeTypeFollowing, targetID)
	return err
}

func (s *SqlSubscriptions) GetFollowing(userID int) ([]int, error) {
	rows, err := s.DB.Query("select to_id from edges where from_id = ? and edge_type = ?", userID, store.EdgeTypeFollowing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []int
	for rows.Next() {
		userID := 0
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}

		result = append(result, userID)
	}

	return result, nil
}

func (s *SqlSubscriptions) CheckFollowing(userID int, targetIds []int) (map[int]bool, error) {
	result := map[int]bool{}

	rows, err := s.DB.Query(fmt.Sprintf("select to_id from edges where from_id = %d and edge_type = %d and to_id in (%s)", userID, store.EdgeTypeFollowing, strings.Join(utils.IdsToStrings(targetIds), ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		targetID := 0
		err = rows.Scan(&targetID)
		if err != nil {
			return nil, err
		}

		result[targetID] = true
	}

	return result, err
}

type Subscriptions interface {
	Follow(userID, targetID int) error
	Unfollow(userID, targetID int) error
	GetFollowing(userID int) ([]int, error)
	CheckFollowing(userID int, targetIds []int) (map[int]bool, error)
}

type MockSubscriptions struct {
	following map[int][]int
}

func (m *MockSubscriptions) Follow(userID, targetID int) error {
	for _, userID := range m.following[userID] {
		if userID == targetID {
			return nil
		}
	}

	m.following[userID] = append(m.following[userID], targetID)
	return nil
}

func (m *MockSubscriptions) Unfollow(userID, targetID int) error {
	var newList []int
	for _, userID := range m.following[userID] {
		if userID != targetID {
			newList = append(newList, userID)
		}
	}
	m.following[userID] = newList

	return nil
}

func (m *MockSubscriptions) GetFollowing(userID int) ([]int, error) {
	return m.following[userID], nil
}

func (m *MockSubscriptions) CheckFollowing(userID int, targetIds []int) (map[int]bool, error) {
	result := map[int]bool{}
	for _, targetID := range m.following[userID] {
		for _, neededID := range targetIds {
			if targetID == neededID {
				result[targetID] = true
			}
		}
	}

	return result, nil
}
