package db

import (
	"context"
	"testing"
	"time"

	"github.com/kratos69/movie-app/util"
	"github.com/stretchr/testify/require"
)

func createRandomMovie(t *testing.T) Movie {
	arg := CreateMovieParams{
		Title:       util.RandomTitle(),
		Description: util.RandomDescription(),
		PosterUrl:   util.RandomPosterURL(),
		GenreID:     util.RandomGenreID(),
	}
	movie, err := testStore.CreateMovie(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, movie)

	require.Equal(t, arg.Title, movie.Title)
	require.Equal(t, arg.Description, movie.Description)
	require.Equal(t, arg.PosterUrl, movie.PosterUrl)
	require.Equal(t, arg.GenreID, movie.GenreID)

	require.NotZero(t, movie.MovieID)
	require.NotZero(t, movie.CreatedAt)

	return movie
}

func TestCreateMovie(t *testing.T) {
	createRandomMovie(t)
}

func TestGetMovie(t *testing.T) {
	movie1 := createRandomMovie(t)
	movie2, err := testStore.GetMovie(context.Background(), movie1.MovieID)
	require.NoError(t, err)
	require.NotEmpty(t, movie2)

	require.Equal(t, movie1.Title, movie2.Title)
	require.Equal(t, movie1.Description, movie2.Description)
	require.Equal(t, movie1.PosterUrl, movie2.PosterUrl)
	require.Equal(t, movie1.GenreID, movie2.GenreID)
	require.WithinDuration(t, movie1.CreatedAt, movie2.CreatedAt, time.Second)
}

func TestDeleteMovie(t *testing.T) {
	movie1 := createRandomMovie(t)
	err := testStore.DeleteMovie(context.Background(), movie1.MovieID)
	require.NoError(t, err)

	movie2, err := testStore.GetMovie(context.Background(), movie1.MovieID)
	require.Error(t, err)
	require.Empty(t, movie2)
}

func TestListMovies(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomMovie(t)
	}

	arg := ListMoviesParams{
		Limit: 5,
		Offset: 5,
	}
	movies, err := testStore.ListMovies(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, movies)
}

func TestUpdateMovie(t *testing.T) {
	movie1 := createRandomMovie(t)

	arg := UpdateMovieParams{
		MovieID:     movie1.MovieID,
		Title:       util.RandomTitle(),
		Description: util.RandomDescription(),
		PosterUrl:   util.RandomPosterURL(),
		GenreID:     util.RandomGenreID(),
	}

	movie2, err := testStore.UpdateMovie(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, movie2)

	require.Equal(t, movie1.MovieID, movie2.MovieID)
}
