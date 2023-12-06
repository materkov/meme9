package nodes

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestStore_Get(t *testing.T) {
	s := Store{
		waiting: map[int]func() int{},
	}

	//objectID, err := s.Add(546512)
	//require.NoError(t, err)
	log.Printf("Start")

	log.Printf("Requesting 1")
	objectID := 105
	objectData1 := s.Get(objectID, nil)
	require.Equal(t, objectData1, 900)
	log.Printf("Got 1")

	objectData2 := s.Get(objectID, nil)
	require.Equal(t, objectData2, 900)
	log.Printf("Got 2")

	objectData3 := s.Get(objectID, nil)
	require.Equal(t, objectData3, 900)
	log.Printf("Got 3")
}
