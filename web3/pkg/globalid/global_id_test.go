package globalid

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGlobalID(t *testing.T) {
	id := Create(PostID{PostID: 15})
	require.Equal(t, "UG9zdElEOnsicG9zdElkIjoxNX0", id) // PostID:{"postId":15}

	objectID, err := Parse(id)
	require.NoError(t, err)
	require.Equal(t, 15, objectID.(*PostID).PostID)
}
