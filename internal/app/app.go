package app

import (
	"errors"
	"log/slog"

	"github.com/findacocktail/backend/internal/pkg/recipes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	ErrGetEvents = errors.New("could not retrieve events")
)

type Server struct {
	port          string
	logger        *slog.Logger
	router        *gin.Engine
	recipeService recipes.Service
}

type API interface {
	Serve()
}

func New(port string, logger *slog.Logger, recipeService recipes.Service) *Server {
	return &Server{
		port:          port,
		logger:        logger,
		recipeService: recipeService,
	}
}

func (s *Server) Serve() {
	router := gin.Default()
	s.router = router

	s.setupConfig(router)
	s.setupEndpoints()

	router.Run("0.0.0.0:" + s.port)
}

func (s *Server) setupConfig(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
}
