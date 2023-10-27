package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFetchURL(t *testing.T) {
	data, err := FetchURL("https://www.kommersant.ru/doc/6309288?from=top_main_1")
	require.NoError(t, err)
	require.NotNil(t, data)
}
