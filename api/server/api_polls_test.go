package server

import (
	"context"
	"github.com/materkov/meme9/api/pb/github.com/materkov/meme9/api"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPI_PollsAdd(t *testing.T) {
	srv := PollServer{}
	closer := createTestDB(t)
	defer closer()

	ctx := context.WithValue(context.Background(), CtxViewerKey, &Viewer{UserID: 15})

	resp, err := srv.Add(ctx, &api.PollsAddReq{
		Question: "my question",
		Answers:  []string{"answer 1", "answer 2"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Id)

	_, err = srv.Vote(ctx, &api.PollsVoteReq{
		PollId:    resp.Id,
		AnswerIds: []string{resp.Answers[0].Id},
	})
	require.NoError(t, err)

	listResp, err := srv.List(ctx, &api.PollsListReq{
		Ids: []string{resp.Id},
	})
	require.NoError(t, err)
	require.Len(t, listResp.Items, 1)
	require.Equal(t, int32(1), listResp.Items[0].Answers[0].VotedCount)
	require.True(t, listResp.Items[0].Answers[0].IsVoted)
}
