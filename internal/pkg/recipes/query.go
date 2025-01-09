package recipes

import (
	"errors"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"

	"github.com/samber/lo"
)

var (
	ErrNotFound = errors.New("cocktail not found")
)

// Search return a list of matches given terms
func (r *service) Search(includedTerms []string, notIncludedTerms []string) ([]Recipe, error) {
	buildQuery := lo.Map(includedTerms, func(term string, _ int) query.Query {
		return bleve.NewQueryStringQuery(term)
	})

	buildQuery = append(buildQuery, lo.Map(notIncludedTerms, func(term string, _ int) query.Query {
		return bleve.NewQueryStringQuery("-" + term)
	})...)

	query := bleve.NewConjunctionQuery(
		buildQuery...,
	)

	searchResult, err := r.recipesIndex.Search(bleve.NewSearchRequest(query))
	if err != nil {
		return nil, err
	}

	if len(searchResult.Hits) == 0 {
		return nil, ErrNotFound
	}

	var results []Recipe
	for _, hit := range searchResult.Hits {
		results = append(results, r.recipesMap[hit.ID])
	}

	return results, nil
}

// RecipeByName returns a recipe by a specific name
func (r *service) RecipeByName(name string) (*Recipe, error) {
	recipe, found := r.recipesMap[name]
	if !found {
		return nil, ErrNotFound
	}

	return &recipe, nil
}
