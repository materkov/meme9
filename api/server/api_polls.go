package server

import (
	"context"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/materkov/meme9/api/src/pkg"
	"github.com/materkov/meme9/api/src/pkg/tracer"
	"github.com/materkov/meme9/api/src/pkg/utils"
	"github.com/materkov/meme9/api/src/store"
	"github.com/materkov/meme9/api/src/store2"
	"github.com/twitchtv/twirp"
	"strconv"
)

type PollServer struct {
}

func transformPollsMany(ctx context.Context, polls []*store.Poll, viewerID int) []*api.Poll {
	defer tracer.FromCtx(ctx).StartChild("transformPollsMany").Stop()

	result := make([]*api.Poll, len(polls))

	var answerIds []int
	for _, poll := range polls {
		answerIds = append(answerIds, poll.AnswerIds...)
	}

	answers, err := store2.GlobalStore.PollAnswers.Get(answerIds)
	pkg.LogErr(err)

	counters, isVoted, err := store2.GlobalStore.Votes.LoadAnswersMany(ctx, answerIds, viewerID)
	pkg.LogErr(err)

	for i, poll := range polls {
		pollAnswers := make([]*api.PollAnswer, len(poll.AnswerIds))
		for i, answerID := range poll.AnswerIds {
			answer, ok := answers[answerID]
			if !ok {
				pollAnswers[i] = &api.PollAnswer{
					Id: strconv.Itoa(answerID),
				}
				continue
			}

			pollAnswers[i] = &api.PollAnswer{
				Id:         strconv.Itoa(answer.ID),
				Answer:     answer.Answer,
				VotedCount: int32(counters[answerID]),
				IsVoted:    isVoted[answerID],
			}
		}

		result[i] = &api.Poll{
			Id:       strconv.Itoa(poll.ID),
			Question: poll.Question,
			Answers:  pollAnswers,
		}
	}

	return result
}

func (p *PollServer) Add(ctx context.Context, r *api.PollsAddReq) (*api.Poll, error) {
	var err error

	viewer := ctx.Value(CtxViewerKey).(*Viewer)

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

func (p *PollServer) Vote(ctx context.Context, r *api.PollsVoteReq) (*api.Void, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	pollID, _ := strconv.Atoi(r.PollId)
	polls, err := store2.GlobalStore.Polls.Get([]int{pollID})
	if err != nil {
		return nil, err
	} else if polls[pollID] == nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "PollNotFound")
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

	return &api.Void{}, nil
}

func (*PollServer) DeleteVote(ctx context.Context, r *api.PollsDeleteVoteReq) (*api.Void, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	pollID, _ := strconv.Atoi(r.PollId)
	polls, err := store2.GlobalStore.Polls.Get([]int{pollID})
	if err != nil {
		return nil, err
	} else if polls[pollID] == nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "PollNotFound")
	}

	err = store2.GlobalStore.Votes.RemoveVote(viewer.UserID, polls[pollID].AnswerIds)
	pkg.LogErr(err)

	return &api.Void{}, err
}

// TODO need this method?
func (*PollServer) List(ctx context.Context, r *api.PollsListReq) (*api.PollsList, error) {
	viewer := ctx.Value(CtxViewerKey).(*Viewer)

	polls, err := store2.GlobalStore.Polls.Get(utils.IdsToInts(r.Ids))
	if err != nil {
		return nil, err
	}

	var pollsList []*store.Poll
	for _, poll := range polls {
		pollsList = append(pollsList, poll)
	}

	return &api.PollsList{Items: transformPollsMany(ctx, pollsList, viewer.UserID)}, nil
}
