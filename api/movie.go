package api

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/kratos69/movie-app/db/sqlc"
	"github.com/kratos69/movie-app/util"
)

// to create a movie in database
func (server *Server) createMovie(ctx *gin.Context) {
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	strGenreID := ctx.PostForm("genre_id")

	if title == "" || description == "" || strGenreID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "title, description & genre_id are required",
		})
		return
	}

	genreID, err := strconv.Atoi(strGenreID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "genre_id must be an integer"})
		return
	}

	// (Optional) Check if genre exists
	// _, err = server.store.gen(ctx, int32(genreID))
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "genre not found"})
	// 	return
	// }

	posterURL := ""
	fileHeader, _ := ctx.FormFile("poster_url")
	if fileHeader != nil {
		posterURL, err = uploadToCloud(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
	}

	arg := db.CreateMovieParams{
		Title:       title,
		Description: description,
		PosterUrl:   posterURL,
		GenreID:     int32(genreID),
	}

	movie, err := server.store.CreateMovie(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "movie created successfully",
		"data":    movie,
	})
}


// Page 1 (first 50 movies): LIMIT 50 OFFSET 0
// Page 2 (next 50):→ LIMIT 50 OFFSET 50
// Page 3 (next 50):→ LIMIT 50 OFFSET 100
// Offset = limit * (page - 1)
func (server *Server) listAllMovies(ctx *gin.Context) {
	// get query params (e.g., ?page=1&limit=50)
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	offset := (page - 1) * limit

	arg := db.ListMoviesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	movies, err := server.store.ListMovies(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "could not fetch movies"})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

type movieIDStruct struct {
	MovieID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getMovieByID(ctx *gin.Context) {
	var req movieIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	movie, err := server.store.GetMovie(ctx, req.MovieID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,
				gin.H{"error": "movie not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

// update a movie
func (server *Server) updateMovie(ctx *gin.Context) {
	var req movieIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	movie, err := server.store.GetMovie(ctx, req.MovieID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,
				gin.H{"error": "movie not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	strGenreID := ctx.PostForm("genre_id")

	if title == "" {
		title = movie.Title
	}

	if description == "" {
		description = movie.Description
	}

	var genreID int32
	if strGenreID == "" {
		genreID = movie.GenreID // fallback to existing
	} else {
		// convert to int
		parsedID, err := strconv.Atoi(strGenreID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid genre id",
			})
			return
		}
		genreID = int32(parsedID)
	}

	// Optional poster upload
	_, err = ctx.FormFile("poster_url")
	var posterURL string
	if err == nil {
		posterURL, err = uploadToCloud(ctx)
		if err != nil {
			// uploadToCloud already handles JSON error response
			return
		}
	} else {
		// No file uploaded → use old URL
		posterURL = movie.PosterUrl
	}

	arg := db.UpdateMovieParams{
		MovieID:     movie.MovieID,
		Title:       title,
		Description: description,
		PosterUrl:   posterURL,
		GenreID:     int32(genreID),
	}

	updatedMovie, err := server.store.UpdateMovie(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedMovie)
}

// delete a movie
func (server *Server) deleteMovie(ctx *gin.Context) {
	var req movieIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	movie, err := server.store.GetMovie(ctx, req.MovieID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound,
				gin.H{"error": "movie not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	cloudService, err := util.NewCloudinaryService()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	// Delete image from Cloudinary
	// extract cloud public ID from URL
	if movie.PosterUrl != "" {
		publicID := extractPublicID(movie.PosterUrl)
		// fmt.Printf("Extracted Public ID: %s\n", publicID)

		if err := cloudService.DeleteImage(ctx, publicID); err != nil {
			log.Printf("Failed to delete image from Cloudinary: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
	}

	err = server.store.DeleteMovie(ctx, movie.MovieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,
		gin.H{"message": "movie deleted"})
}

// ------------------------------------------------------------------//
// ------------------------------Helper Funcs------------------------//
// ------------------------------------------------------------------//

func uploadToCloud(ctx *gin.Context) (string, error) {
	file, err := ctx.FormFile("poster_url")
	if err != nil {
		// checking file size (less than 5 mb)
		if file.Size > 5<<20 {
			ctx.JSON(http.StatusBadRequest, errResponse(err))
			return "", err
		}

		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return "", err
	}

	// check format validity
	if !containsValidFormat(file.Header.Get("Content-Type")) {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return "", err
	}

	// Upload the image locally
	err = ctx.SaveUploadedFile(file, "uploads/"+file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return "", err
	}

	cloudService, err := util.NewCloudinaryService()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return "", err
	}

	imageUrl, err := cloudService.UploadImage(ctx, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return "", err
	}

	return imageUrl, nil
}

// Helper function to check if a string exists in a slice
func containsValidFormat(item string) bool {
	slice := []string{"image/png", "image/jpeg", "image/jpg", "image/gif"}

	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to extract public ID from a Cloudinary URL
func extractPublicID(url string) string {
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]
	publicID := strings.TrimSuffix(lastPart, filepath.Ext(lastPart)) // Remove file extension

	// Extract folder path if it exists
	if len(parts) > 7 { // Cloudinary path structure
		folderPath := strings.Join(parts[7:len(parts)-1], "/") // Preserve folder structure
		publicID = folderPath + "/" + publicID
	}

	return publicID
}
