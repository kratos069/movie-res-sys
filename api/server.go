package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/kratos69/movie-app/db/sqlc"
	"github.com/kratos69/movie-app/token"
	"github.com/kratos69/movie-app/util"
)

// servers HTTP requests for the insta-app
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// Creates HTTP server and Setup Routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Routes
	server.setupRoutes()

	return server, nil
}

func (server *Server) setupRoutes() {
	router := gin.Default()

	// cors
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://127.0.0.1:5173", "http://localhost:5173"}, // frontend origin
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	AllowCredentials: true,
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	MaxAge:           12 * time.Hour,
	// }))

	// routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	router.GET("/movies", server.listAllMovies)
	router.GET("/movies/:id", server.getMovieByID)

	// for both users and admins
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker,
		[]string{util.AdminRole, util.CustomerRole}))
	authRoutes.GET("/users/:user_id", server.getUserByID)

	// for only admins
	adminRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker,
		[]string{util.AdminRole}))
	adminRoutes.POST("/movies", server.createMovie)
	adminRoutes.PUT("/movies/:id", server.updateMovie)
	adminRoutes.DELETE("/movies/:id", server.deleteMovie)

	server.router = router

}

// Starts and runs HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
