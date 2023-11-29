package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPI_PollsAdd(t *testing.T) {
	api := API{}
	closer := createTestDB(t)
	defer closer()

	resp, err := api.PollsAdd(&Viewer{UserID: 15}, &PollsAddReq{
		Question: "my question",
		Answers:  []string{"answer 1", "answer 2"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.ID)

	_, err = api.PollsVote(&Viewer{UserID: 15}, &PollsVoteReq{
		PollID:    resp.ID,
		AnswerIds: []string{resp.Answers[0].ID},
	})
	require.NoError(t, err)

	listResp, err := api.PollsList(&Viewer{UserID: 15}, &PollsListReq{
		Ids: []string{resp.ID},
	})
	require.NoError(t, err)
	require.Len(t, listResp, 1)
	require.Equal(t, 1, listResp[0].Answers[0].VotedCount)
	require.True(t, listResp[0].Answers[0].IsVoted)
}
