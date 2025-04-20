package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kratos69/movie-app/util"
	"github.com/stretchr/testify/require"
)

func createRandomShowtime(t *testing.T) Showtime {
	// First create a movie since showtime depends on it
	movie := createRandomMovie(t)

	// Create timestamp in pgtype format
	now := time.Now()
	startTime := pgtype.Timestamp{}
	err := startTime.Scan(now)
	require.NoError(t, err)

	arg := CreateShowtimeParams{
		MovieID:   movie.MovieID,
		StartTime: startTime,
		Price:     util.RandomPrice(),
	}

	showtime, err := testStore.CreateShowtime(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, showtime)

	require.Equal(t, arg.MovieID, showtime.MovieID)
	require.Equal(t, arg.Price, showtime.Price)

	require.NotZero(t, showtime.ShowtimeID)
	require.NotZero(t, showtime.CreatedAt)

	return showtime
}

func TestCreateShowtime(t *testing.T) {
	createRandomShowtime(t)
}

func TestGetShowtime(t *testing.T) {
	showtime1 := createRandomShowtime(t)

	showtime2, err := testStore.GetShowtime(context.Background(), showtime1.ShowtimeID)
	require.NoError(t, err)
	require.NotEmpty(t, showtime2)

	require.Equal(t, showtime1.ShowtimeID, showtime2.ShowtimeID)
	require.Equal(t, showtime1.MovieID, showtime2.MovieID)
	require.Equal(t, showtime1.Price, showtime2.Price)
	require.WithinDuration(t, showtime1.StartTime.Time, showtime2.StartTime.Time, time.Second)
	require.WithinDuration(t, showtime1.CreatedAt, showtime2.CreatedAt, time.Second)
}
