package api

import (
	"context"
	"github.com/materkov/meme9/web6/src/pkg"
	"github.com/materkov/meme9/web6/src/pkg/tracer"
	"github.com/materkov/meme9/web6/src/pkg/utils"
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

	answers, err := store2.GlobalStore.PollAnswers.Get(answerIds)
	pkg.LogErr(err)

	counters, isVoted, err := store2.GlobalStore.Votes.LoadAnswersMany(ctx, answerIds, viewerID)
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

		err = store2.GlobalStore.PollAnswers.Add(&answer)
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
	err = store2.GlobalStore.Polls.Add(poll)
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
	polls, err := store2.GlobalStore.Polls.Get([]int{pollID})
	if err != nil {
		return nil, err
	} else if polls[pollID] == nil {
		return nil, Error("PollNotFound")
	}

	poll := polls[pollID]

	pollAnswers := poll.GetAnswersMap()

	var answerIds []int
	for _, answerIDStr := range r.AnswerIds {
		answerID, _ := strconv.Atoi(answerIDStr)
		if answerID > 0 && pollAnswers[answerID] {
			answerIds = append(answerIds, answerID)
		}
	}

	err = store2.GlobalStore.Votes.Vote(viewer.UserID, answerIds)
	if err != nil {
		return nil, err
	}

	return &Void{}, nil
}

type PollsDeleteVoteReq struct {
	PollID string `json:"pollId"`
}

func (*API) PollsDeleteVote(viewer *Viewer, r *PollsDeleteVoteReq) (*Void, error) {
	pollID, _ := strconv.Atoi(r.PollID)
	polls, err := store2.GlobalStore.Polls.Get([]int{pollID})
	if err != nil {
		return nil, err
	} else if polls[pollID] == nil {
		return nil, Error("PollNotFound")
	}

	err = store2.GlobalStore.Votes.RemoveVote(viewer.UserID, polls[pollID].AnswerIds)
	pkg.LogErr(err)

	return &Void{}, err
}

type PollsListReq struct {
	Ids []string `json:"ids"`
}

func (*API) PollsList(ctx context.Context, viewer *Viewer, r *PollsListReq) ([]*Poll, error) {
	polls, err := store2.GlobalStore.Polls.Get(utils.IdsToInts(r.Ids))
	if err != nil {
		return nil, err
	}

	var pollsList []*store.Poll
	for _, poll := range polls {
		pollsList = append(pollsList, poll)
	}

	return transformPollsMany(ctx, pollsList, viewer.UserID), nil
}
