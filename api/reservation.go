package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kratos69/movie-app/db/sqlc"
	"github.com/kratos69/movie-app/token"
)

type reserveSeatsRequest struct {
	ShowtimeID int32   `json:"showtime_id" binding:"required"`
	SeatIDs    []int32 `json:"seat_ids" binding:"required,min=1"`
}

//   "showtime_id": 12,
//  "seat_ids": [5, 6, 7]
func (server *Server) reserveSeats(ctx *gin.Context) {
	var req reserveSeatsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ReserveMultipleSeatsTxParams{
		UserID:     authPayload.UserID,
		ShowtimeID: req.ShowtimeID,
		SeatIDs:    req.SeatIDs,
	}

	result, err := server.store.ReserveMultipleSeatsTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type cancelReservationRequestUri struct {
	ResID int64 `uri:"id" binding:"required,min=1"` // ReservationID from URL
}

func (server *Server) cancelReservation(ctx *gin.Context) {
	var uri cancelReservationRequestUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid reservation id"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CancelReservationParams{
		ReservationID: uri.ResID,
		UserID:        authPayload.UserID,
	}

	err := server.store.CancelReservationTx(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "reservation cancelled"})
}

func (server *Server) listReservationsByUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	reservations, err := server.store.ListReservationsByUser(
		ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "could not fetch reservations"})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}
