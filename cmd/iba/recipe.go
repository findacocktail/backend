package iba

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/findacocktail/backend/cmd/model"
	"github.com/findacocktail/backend/cmd/parsing"
	"golang.org/x/net/html"
)

func (p *ibaParser) GetRecipe(recipeLink string) (*model.Recipe, error) {
	req, err := http.NewRequest(http.MethodGet, recipeLink, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	token, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	cocktailName, err := parsing.GetNode(token, "h1")
	if err != nil {
		return nil, err
	}

	youtubeLink, err := parsing.GetAttributeStartsWith(token, "href", "https://www.youtube.com/watch")
	if err != nil {
		fmt.Println("could not find link", err)
	}

	recipe := model.Recipe{
		Name:        strings.TrimSpace(cocktailName.FirstChild.Data),
		YoutubeLink: youtubeLink,
	}

	image, err := parsing.GetNodeByAttribute(token, "fetchpriority", "high")
	if err != nil {
		return nil, err
	}

	for _, attr := range image.Attr {
		if attr.Key == "src" {
			recipe.ImageURL = attr.Val
			break
		}
	}

	ingredients, err := parseIngredients(token)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients = ingredients

	method, err := parseListOfP(token, "Method")
	if err != nil {
		return nil, err
	}
	recipe.Method = method

	garnish, err := parseListOfP(token, "Garnish")
	if err != nil {
		return nil, err
	}
	recipe.Garnish = garnish

	return &recipe, nil
}
