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

func (r *service) Search(terms []string) ([]Recipe, error) {
	query := bleve.NewConjunctionQuery(
		lo.Map(terms, func(term string, _ int) query.Query {
			return bleve.NewQueryStringQuery(term)
		})...,
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
