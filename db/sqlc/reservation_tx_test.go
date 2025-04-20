package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReserveMultipleSeatsTx(t *testing.T) {
	user := createRandomUser(t)
	showtime := createRandomShowtime(t)

	// Get random available seats for this showtime
	availableSeats := getRandomAvailableSeats(t, showtime.ShowtimeID, 3)
	seatIDs := []int32{
		availableSeats[0].SeatID,
		availableSeats[1].SeatID,
		availableSeats[2].SeatID,
	}

	arg := ReserveMultipleSeatsTxParams{
		UserID:     user.UserID,
		ShowtimeID: showtime.ShowtimeID,
		SeatIDs:    seatIDs,
	}

	result, err := testStore.ReserveMultipleSeatsTx(context.Background(),
		arg)
	require.NoError(t, err)
	require.Len(t, result.Reservations, len(seatIDs))

	for i, res := range result.Reservations {
		require.Equal(t, user.UserID, res.UserID)
		require.Equal(t, showtime.ShowtimeID, res.ShowtimeID)
		require.Equal(t, seatIDs[i], res.SeatID)
		require.NotZero(t, res.ReservationID)
		require.WithinDuration(t, time.Now(), res.ReservedAt, time.Second)
	}
}

func TestCancelReservationTx(t *testing.T) {
	user := createRandomUser(t)
	showtime := createRandomShowtime(t)

	// Pick a random seat to reserve
	seats, err := testStore.ListAvailableSeatsForShowtime(
		context.Background(), showtime.ShowtimeID)
	require.NoError(t, err)
	require.NotEmpty(t, seats)

	seat := seats[0]

	// Reserve the seat
	reserveArg := ReserveSeatParams{
		UserID:     user.UserID,
		ShowtimeID: showtime.ShowtimeID,
		SeatID:     seat.SeatID,
	}
	reserveResult, err := testStore.ReserveSeat(context.Background(),
		reserveArg)
	require.NoError(t, err)

	// Cancel the reservation
	cancelArg := CancelReservationParams{
		ReservationID: reserveResult.ReservationID,
		UserID:        reserveResult.UserID,
	}
	err = testStore.CancelReservationTx(context.Background(), cancelArg)
	require.NoError(t, err)

	// Try to get the reservation again — it should not exist
	getReservation, err := testStore.ListReservationsByUser(
		context.Background(), cancelArg.UserID)
	require.NoError(t, err)
	require.Len(t, getReservation, 0)
	require.Empty(t, getReservation)
}

func getRandomAvailableSeats(t *testing.T, showtimeID int32, n int) []Seat {
	seats, err := testStore.ListAvailableSeatsForShowtime(
		context.Background(), showtimeID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(seats), n)

	return seats[:n]
}
