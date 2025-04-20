package db

import (
	"context"
	"fmt"
	"strings"
)

// To prevent race conditions, like two users reserving
// the same seat at the same time.

type ReserveMultipleSeatsTxParams struct {
	UserID     int64   `json:"user_id"`
	ShowtimeID int32   `json:"showtime_id"`
	SeatIDs    []int32 `json:"seat_ids"`
}

type ReserveMultipleSeatsTxResult struct {
	Reservations []Reservation `json:"reservations"`
}

// Reserve a seat for a user/users only if it's not already reserved for
// that showtime. This operation must be atomic.
func (store *SQLStore) ReserveMultipleSeatsTx(
	ctx context.Context, arg ReserveMultipleSeatsTxParams,
) (ReserveMultipleSeatsTxResult, error) {
	var result ReserveMultipleSeatsTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// Step 1: get available seats
		availableSeats, err := q.ListAvailableSeatsForShowtime(ctx,
			arg.ShowtimeID)
		if err != nil {
			return err
		}

		// put available seats in map
		availableMap := make(map[int32]bool)
		for _, s := range availableSeats {
			availableMap[s.SeatID] = true
		}

		// Step 2: validate all requested seats are available
		for _, seatID := range arg.SeatIDs {
			if !availableMap[seatID] {
				return fmt.Errorf(
					"seat %d is not available for showtime %d",
					seatID, arg.ShowtimeID)
			}
		}

		// Step 3: insert each seat one by one
		for _, seatID := range arg.SeatIDs {
			res, err := q.ReserveSeat(ctx, ReserveSeatParams{
				UserID:     arg.UserID,
				ShowtimeID: arg.ShowtimeID,
				SeatID:     seatID,
			})
			if err != nil {
				// Handle DB unique constraint (concurrent race case)
				if strings.Contains(err.Error(), "duplicate key") {
					return fmt.Errorf("seat %d already reserved", seatID)
				}
				return err
			}
			result.Reservations = append(result.Reservations, res)
		}

		return nil
	})

	return result, err
}

// Cancelling reservation in tx
func (store *SQLStore) CancelReservationTx(ctx context.Context, 
	arg CancelReservationParams) error {
	return store.execTx(ctx, func(q *Queries) error {
		return q.CancelReservation(ctx, CancelReservationParams{
			ReservationID: arg.ReservationID,
			UserID:        arg.UserID,
		})
	})
}

