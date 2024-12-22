package recipes

import (
	"bytes"
	"embed"
	"encoding/json"
	"log/slog"

	"github.com/blevesearch/bleve"
	"github.com/samber/lo"
)

//go:embed static/20241221.json
var content embed.FS

const (
	recipesFileName string = "static/20241221.json"
)

type Service interface {
	Search(includedTerms []string, notIncludedTerms []string) ([]Recipe, error)
}

type service struct {
	logger       *slog.Logger
	recipesIndex bleve.Index
	recipesMap   map[string]Recipe
}

func New(logger *slog.Logger) (*service, error) {
	recipes, err := parseRecipes()
	if err != nil {
		return nil, err
	}

	index, err := readAndIndexFiles(recipes)
	if err != nil {
		return nil, err
	}

	recipesMap := lo.Associate(recipes, func(r Recipe) (string, Recipe) {
		return r.Name, r
	})

	return &service{
		logger:       logger,
		recipesIndex: *index,
		recipesMap:   recipesMap,
	}, nil
}

func readAndIndexFiles(recipes []Recipe) (*bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	for _, recipe := range recipes {
		err = index.Index(recipe.Name, recipe)
		if err != nil {
			return nil, err
		}
	}

	return &index, nil
}

func parseRecipes() ([]Recipe, error) {
	data, err := content.ReadFile(recipesFileName)
	if err != nil {
		return nil, err
	}

	var recipes []Recipe
	err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&recipes)
	return recipes, err
}
