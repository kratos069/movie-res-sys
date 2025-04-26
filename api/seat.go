package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kratos69/movie-app/db/sqlc"
)

// list all the seats for a specific showtime_id
func (server *Server) listSeatsForShowtime(ctx *gin.Context) {
	var uri struct {
		ID int32 `uri:"id" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "invalid showtime ID"})
		return
	}

	seats, err := server.store.ListSeatsForShowtime(ctx, uri.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "could not fetch seats"})
		return
	}

	grouped := make(map[int32][]db.ListSeatsForShowtimeRow)
	for _, seat := range seats {
		grouped[seat.Row] = append(grouped[seat.Row], seat)
	}

	ctx.JSON(http.StatusOK, grouped)
}
