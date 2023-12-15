package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/store"
	"github.com/materkov/meme9/web6/src/store2"
	"strconv"
)

type PollsAddReq struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

type Poll struct {
	ID       string        `json:"id"`
	Question string        `json:"question"`
	Answers  []*PollAnswer `json:"answers"`
}

type PollAnswer struct {
	ID         string `json:"id"`
	Answer     string `json:"answer"`
	VotedCount int    `json:"voted,omitempty"`
	IsVoted    bool   `json:"isVoted,omitempty"`
}

func transformPollsMany(ctx context.Context, polls []*store.Poll, viewerID int) []*Poll {
	defer tracer.FromCtx(ctx).StartChild("transformPollsMany").Stop()

	result := make([]*Poll, len(polls))

	var answerIds []int
	for _, poll := range polls {
		answerIds = append(answerIds, poll.AnswerIds...)
	}

	answerBytes, err := store.GlobalStore.GetObjectsMany(ctx, answerIds)
	pkg.LogErr(err)

	answers := map[int]*store.PollAnswer{}
	for _, answerID := range answerIds {
		pollAnswer := store.PollAnswer{}
		err = json.Unmarshal(answerBytes[answerID], &pollAnswer)
		pkg.LogErr(err)
		if err == nil {
			pollAnswer.ID = answerID
			answers[pollAnswer.ID] = &pollAnswer
		}
	}

	counters, isVoted, err := store.GlobalStore.LoadAnswersMany(ctx, answerIds, viewerID)
	pkg.LogErr(err)

	for i, poll := range polls {
		pollAnswers := make([]*PollAnswer, len(poll.AnswerIds))
		for i, answerID := range poll.AnswerIds {
			answer, ok := answers[answerID]
			if !ok {
				pollAnswers[i] = &PollAnswer{
					ID: strconv.Itoa(answerID),
				}
				continue
			}

			pollAnswers[i] = &PollAnswer{
				ID:         strconv.Itoa(answer.ID),
				Answer:     answer.Answer,
				VotedCount: counters[answerID],
				IsVoted:    isVoted[answerID],
			}
		}

		result[i] = &Poll{
			ID:       strconv.Itoa(poll.ID),
			Question: poll.Question,
			Answers:  pollAnswers,
		}
	}

	return result
}

func (*API) PollsAdd(ctx context.Context, viewer *Viewer, r *PollsAddReq) (*Poll, error) {
	var err error

	answerIds := make([]int, len(r.Answers))
	for i, inputAnswer := range r.Answers {
		answer := store.PollAnswer{
			Answer: inputAnswer,
		}

		answer.ID, err = store2.GlobalStore.Nodes.Add(store.ObjTypePollAnswer, &answer)
		if err != nil {
			return nil, err
		}

		answerIds[i] = answer.ID
	}

	poll := &store.Poll{
		ID:        0,
		UserID:    viewer.UserID,
		Question:  r.Question,
		AnswerIds: answerIds,
	}
	poll.ID, err = store2.GlobalStore.Nodes.Add(store.ObjTypePoll, poll)
	if err != nil {
		return nil, err
	}

	return transformPollsMany(ctx, []*store.Poll{poll}, viewer.UserID)[0], nil
}

type PollsVoteReq struct {
	PollID    string   `json:"pollId"`
	AnswerIds []string `json:"answerIds"`
}

func (*API) PollsVote(viewer *Viewer, r *PollsVoteReq) (*Void, error) {
	pollID, _ := strconv.Atoi(r.PollID)
	poll, err := store.GetPoll(pollID)
	if errors.Is(err, store.ErrObjectNotFound) {
		return nil, Error("PollNotFound")
	} else if err != nil {
		return nil, err
	}

	alreadyVoted := false
	for _, answerID := range poll.AnswerIds {
		_, err := store.GlobalStore.GetEdge(answerID, viewer.UserID, store.EdgeTypeVoted)
		if err == nil {
			alreadyVoted = true
			break
		} else if !errors.Is(err, store.ErrNoEdge) {
			return nil, err
		}
	}

	if !alreadyVoted {
		for _, answerIdStr := range r.AnswerIds {
			answerID, _ := strconv.Atoi(answerIdStr)
			if answerID <= 0 {
				continue
			}

			err = store.GlobalStore.AddEdge(answerID, viewer.UserID, store.EdgeTypeVoted)
			pkg.LogErr(err)
		}
	}

	return &Void{}, nil
}

type PollsDeleteVoteReq struct {
	PollID string `json:"pollId"`
}

func (*API) PollsDeleteVote(viewer *Viewer, r *PollsDeleteVoteReq) (*Void, error) {
	pollID, _ := strconv.Atoi(r.PollID)
	poll, err := store.GetPoll(pollID)
	if errors.Is(err, store.ErrObjectNotFound) {
		return nil, Error("PollNotFound")
	} else if err != nil {
		return nil, err
	}

	for _, answerID := range poll.AnswerIds {
		err := store.GlobalStore.DelEdge(answerID, viewer.UserID, store.EdgeTypeVoted)
		pkg.LogErr(err)
	}
	return &Void{}, err
}

type PollsListReq struct {
	Ids []string `json:"ids"`
}

func (*API) PollsList(ctx context.Context, viewer *Viewer, r *PollsListReq) ([]*Poll, error) {
	pollID, _ := strconv.Atoi(r.Ids[0])
	poll, err := store.GetPoll(pollID)
	if err != nil {
		return nil, err
	}

	return transformPollsMany(ctx, []*store.Poll{poll}, viewer.UserID), nil
}
