package liquorcom

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/findacocktail/backend/cmd/model"
	"github.com/findacocktail/backend/cmd/parsing"
	archiveorg "github.com/findacocktail/backend/internal/pkg/archive.org"
	"golang.org/x/net/html"
)

func (p *liquorParser) GetRecipe(recipeLink string) (*model.Recipe, error) {
	if p.cache {
		archiveService := archiveorg.New()
		newLink, err := archiveService.GetLastSnapshot(recipeLink)
		if err != nil {
			return nil, err
		}
		recipeLink = newLink
	}

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

	splitString := strings.Split(recipeLink, "https://www.liquor.com")
	waybackURL := splitString[0]

	waybackURL = waybackURL[:len(waybackURL)-1]
	waybackURL += "im_/"

	imageLink, err := parsing.GetNodeByAttribute(token, "property", "og:image")
	if err != nil {
		slog.Error("could not find link", slog.Any("err", err), slog.Any("link", recipeLink))
	} else {
		for _, attr := range imageLink.Attr {
			if attr.Key == "content" {
				recipe.ImageURL = strings.TrimPrefix(attr.Val, waybackURL)
			}
		}
	}

	ingredients, err := parseIngredients(token)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients = ingredients

	method, err := parseMethod(token)
	if err != nil {
		return nil, err
	}
	recipe.Method = method

	return &recipe, nil
}
