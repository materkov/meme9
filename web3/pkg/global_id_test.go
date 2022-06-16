package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGlobalID(t *testing.T) {
	id := GetGlobalID(GlobalIDPost, 15)
	require.Equal(t, "UG9zdDoxNQ", id)

	objectType, objectID, err := ParseGlobalID(id)
	require.NoError(t, err)
	require.Equal(t, GlobalIDPost, objectType)
	require.Equal(t, 15, objectID)
}
