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
		ingredient := model.Ingredient{}

		quantityNode, err := parsing.GetNodeIfAttributeExists(node, "data-ingredient-quantity")
		if err == nil {
			quantity, err := strconv.ParseFloat(quantityNode.Data, 64)
			if err == nil {
				ingredient.Amount = quantity
			}
		}

		unitNode, err := parsing.GetNodeIfAttributeExists(node, "data-ingredient-unit")
		if err == nil {
			unit, err := strconv.ParseFloat(unitNode.Data, 64)
			if err == nil {
				ingredient.Amount = unit
			}
		}

		nameNode, err := parsing.GetNodeIfAttributeExists(node, "data-ingredient-name")
		if err == nil {
			ingredient.Description = nameNode.Data
		}

		ingredients = append(ingredients, &ingredient)
	}

	return ingredients, nil
}
