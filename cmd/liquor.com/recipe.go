package liquorcom

import (
	"net/http"
	"strings"

	"github.com/findacocktail/backend/cmd/model"
	"github.com/findacocktail/backend/cmd/parsing"
	"golang.org/x/net/html"
)

func (p *liquorParser) GetRecipe(recipeLink string) (*model.Recipe, error) {
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

	cocktailName, err := parsing.GetNodeByAttribute(token, "class", "heading__title")
	if err != nil {
		return nil, err
	}

	recipe := model.Recipe{
		Name: strings.TrimSpace(cocktailName.FirstChild.Data),
	}

	image, err := parsing.GetNodeByAttribute(token, "id", "mntl-sc-block-image_1-0")
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

	/*
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
	*/
	return &recipe, nil
}
