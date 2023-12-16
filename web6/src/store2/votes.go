package store2

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/pkg/utils"
	"github.com/materkov/meme9/web6/src/store"
	"strings"
)

type SqlVotes struct {
	DB *sql.DB
}

func (s *SqlVotes) Vote(userID int, answerIds []int) error {
	parts := make([]string, len(answerIds))
	for i, answerID := range answerIds {
		parts[i] = fmt.Sprintf("(%d, %d, %d, unix_timestamp())", answerID, userID, store.EdgeTypeVoted)
	}
	_, err := s.DB.Exec("insert into edges (from_id, to_id, edge_type, date) values " + strings.Join(parts, ","))
	return err
}

func (s *SqlVotes) RemoveVote(userID int, answerIds []int) error {
	parts := make([]string, len(answerIds))
	for i, answerID := range answerIds {
		parts[i] = fmt.Sprintf("(from_id = %d and edge_type = %d and to_id = %d)", answerID, store.EdgeTypeVoted, userID)
	}
	_, err := s.DB.Exec("delete from edges where " + strings.Join(parts, " or "))
	return err
}

func (s *SqlVotes) LoadAnswersMany(ctx context.Context, answerIds []int, viewerID int) (counters map[int]int, isVoted map[int]bool, err error) {
	defer tracer.FromCtx(ctx).StartChild("LoadAnswersMany").Stop()

	if len(answerIds) == 0 {
		return map[int]int{}, map[int]bool{}, nil
	}

	query := `
select from_id, count(*), sum(to_id = %d)
from edges
where from_id in (%s) and edge_type=%d
group by from_id
`
	rows, err := s.DB.Query(fmt.Sprintf(query, viewerID, strings.Join(utils.IdsToStrings(answerIds), ","), store.EdgeTypeVoted))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	isVoted = map[int]bool{}
	counters = map[int]int{}

	for rows.Next() {
		postID, count, isLikedInt := 0, 0, 0
		err = rows.Scan(&postID, &count, &isLikedInt)
		if err != nil {
			return nil, nil, err
		}

		counters[postID] = count
		if isLikedInt > 0 {
			isVoted[postID] = true
		}
	}

	return counters, isVoted, nil
}

type Votes interface {
	Vote(userID int, answerIds []int) error
	RemoveVote(userID int, answerIds []int) error
	LoadAnswersMany(ctx context.Context, answerIds []int, viewerID int) (counters map[int]int, isVoted map[int]bool, err error)
}

type MockVotes struct {
	votes map[int][]int
}

func (m *MockVotes) Vote(userID int, answerIds []int) error {
	for _, answerID := range answerIds {
		found := false
		for _, curUserID := range m.votes[answerID] {
			if userID == curUserID {
				found = true
			}
		}
		if !found {
			m.votes[answerID] = append(m.votes[answerID], userID)
		}
	}
	return nil
}

func (m *MockVotes) RemoveVote(userID int, answerIds []int) error {
	for _, answerID := range answerIds {
		var newList []int
		for _, curUserID := range m.votes[answerID] {
			if userID != curUserID {
				newList = append(newList, curUserID)
			}
		}

		m.votes[answerID] = newList
	}
	return nil
}

func (m *MockVotes) LoadAnswersMany(ctx context.Context, answerIds []int, viewerID int) (counters map[int]int, isVoted map[int]bool, err error) {
	counters = map[int]int{}
	isVoted = map[int]bool{}

	for _, answerID := range answerIds {
		counters[answerID] = len(m.votes[answerID])
		for _, userID := range m.votes[answerID] {
			if userID == viewerID {
				isVoted[answerID] = true
			}
		}
	}

	return counters, isVoted, nil
}
