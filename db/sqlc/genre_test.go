package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListGenre(t *testing.T) {
	genres, err := testStore.ListGenres(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, genres)
	require.Equal(t, len(genres), 10)
}
