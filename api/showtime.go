package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/kratos69/movie-app/db/sqlc"
)

type req struct {
	MovieID   int32  `json:"movie_id" binding:"required,min=1"`
	StartTime string `json:"start_time" binding:"required"` // Format: "2025-05-01T20:00"
	Price     string `json:"price" binding:"required"`      // Format: "9.99"
}

func (server *Server) createShowtime(ctx *gin.Context) {
	var showtimeReq req

	if err := ctx.ShouldBindJSON(&showtimeReq); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid request format"})
		return
	}

	// Parse and validate start_time
	t, err := time.Parse("2006-01-02T15:04", showtimeReq.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid start_time format. Use YYYY-MM-DDTHH:MM"})
		return
	}
	if t.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "start_time cannot be in the past"})
		return
	}
	startTime := pgtype.Timestamp{}
	if err := startTime.Scan(t); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "failed to set start_time"})
		return
	}

	// Parse and validate price
	priceFloat, err := strconv.ParseFloat(showtimeReq.Price, 64)
	if err != nil || priceFloat <= 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "price must be a positive number"})
		return
	}
	price := pgtype.Numeric{}
	if err := price.Scan(showtimeReq.Price); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid price format"})
		return
	}

	arg := db.CreateShowtimeParams{
		MovieID:   showtimeReq.MovieID,
		StartTime: startTime,
		Price:     price,
	}

	showtime, err := server.store.CreateShowtime(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, showtime)
}

func (server *Server) getShowtime(ctx *gin.Context) {
	var uri struct {
		ID int32 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid showtime ID"})
		return
	}

	showtime, err := server.store.GetShowtime(ctx, uri.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{"error": "showtime not found"})
		return
	}

	ctx.JSON(http.StatusOK, showtime)
}

// /showtimes
// /showtimes?date=2025-05-01
func (server *Server) listShowtimes(ctx *gin.Context) {
	dateQuery := ctx.Query("date")

	var start, end time.Time
	var err error

	if dateQuery == "" {
		// No date param = return upcoming showtimes
		start = time.Now().UTC()
		end = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC) // Far future to fetch all
	} else {
		// Parse date format: "2006-01-02"
		start, err = time.Parse("2006-01-02", dateQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		// End of day (23:59:59)
		end = start.Add(24 * time.Hour)
	}

	// Convert to pgtype.Timestamp
	startPg := pgtype.Timestamp{Time: start, Valid: true}
	endPg := pgtype.Timestamp{Time: end, Valid: true}

	// Call modified DB function
	showtimes, err := server.store.ListShowtimesBetween(ctx,
		db.ListShowtimesBetweenParams{
			StartTime:   startPg,
			StartTime_2: endPg,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, showtimes)
}

// delete a showTime
func (server *Server) deleteShowtime(ctx *gin.Context) {
	var uri struct {
		ID int32 `uri:"id" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid showtime ID"})
		return
	}

	err := server.store.DeleteShowtime(ctx, uri.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "unable to delete a showtime"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "showtime deleted",
	})
}
