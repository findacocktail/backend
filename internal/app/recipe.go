package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	IDParam = "id"

	FieldParam = "field"
	ValueParam = "value"
)

func (s *Server) setupEndpoints() {
	usersEndpoint := s.router.Group("/cocktails")
	usersEndpoint.GET("", s.getCocktails)
}

func (s *Server) getCocktails(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
