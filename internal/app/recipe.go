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
	usersEndpoint := s.router.Group("/users")
	usersEndpoint.GET("", s.getUsers)
}

func (s *Server) getUsers(c *gin.Context) {

	c.JSON(http.StatusOK, nil)
}
