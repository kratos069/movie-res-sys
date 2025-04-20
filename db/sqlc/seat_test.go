package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListSeats(t *testing.T) {
	seats, err := testStore.ListAllSeats(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, seats)

	require.Equal(t, len(seats), 50)
}
