package app

import (
	"errors"
	"net/http"

	"github.com/findacocktail/backend/internal/pkg/recipes"
	"github.com/gin-gonic/gin"
)

const (
	IDParam = "id"

	FieldParam = "field"
	ValueParam = "value"
)

func (s *Server) setupEndpoints() {
	cocktailsEndpoint := s.router.Group("/cocktails")
	cocktailsEndpoint.GET("", s.getCocktails)
	cocktailsEndpoint.GET(":name", s.getCocktail)
}

func (s *Server) getCocktails(c *gin.Context) {
	results, err := s.recipeService.Search(
		c.QueryArray("term"),
		c.QueryArray("notIncluded"),
	)
	if errors.Is(err, recipes.ErrNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func (s *Server) getCocktail(c *gin.Context) {
	result, err := s.recipeService.RecipeByName(
		c.Param("name"),
	)
	if errors.Is(err, recipes.ErrNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
