package liquorcom

import (
	"strconv"

	"github.com/findacocktail/backend/cmd/model"
	"github.com/findacocktail/backend/cmd/parsing"
	"golang.org/x/net/html"
)

func parseIngredients(token *html.Node) ([]*model.Ingredient, error) {
	ingredientList, err := parsing.GetNodeByAttribute(token, "class", "structured-ingredients__list text-passage")
	if err != nil {
		return nil, err
	}

	ingredients := []*model.Ingredient{}
	for node := range ingredientList.ChildNodes() {
		if node.Data != "li" {
			continue
		}

		ingredient := model.Ingredient{}
		quantityNode, err := parsing.GetNodeByAttribute(node, "data-ingredient-quantity", "true")
		if err == nil && quantityNode.FirstChild != nil {
			quantity, err := strconv.ParseFloat(quantityNode.FirstChild.Data, 64)
			if err == nil {
				ingredient.Amount = quantity
			}
		}

		unitNode, err := parsing.GetNodeByAttribute(node, "data-ingredient-unit", "true")
		if err == nil && unitNode.FirstChild != nil {
			ingredient.Scale = unitNode.FirstChild.Data

		}

		nameNode, err := parsing.GetNodeByAttribute(node, "data-ingredient-name", "true")
		if err == nil && nameNode.FirstChild != nil {
			ingredient.Description = nameNode.FirstChild.Data
		}

		ingredients = append(ingredients, &ingredient)
	}

	return ingredients, nil
}
