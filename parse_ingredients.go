package main

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func parseIngredients(root *html.Node) ([]*Ingredient, error) {
	node, err := getNode(root, "Ingredients")
	if err != nil {
		return nil, err
	}

	topParent := node.Parent.Parent.Parent
	nextNode := topParent.NextSibling.NextSibling

	ulList, err := getNode(nextNode, "ul")
	if err != nil {
		return nil, err
	}

	child := ulList.FirstChild
	ingredients := []*Ingredient{}
	for child != nil {
		if child.Data == "li" {
			ingredientSplit := strings.Split(child.FirstChild.Data, " ")
			amount, err := strconv.ParseFloat(ingredientSplit[0], 64)
			if err != nil {
				return nil, err
			}

			ingredients = append(ingredients, &Ingredient{
				Amount:      amount,
				Scale:       ingredientSplit[1],
				Description: strings.Join(ingredientSplit[2:], " "),
			})
		}
		child = child.NextSibling
	}

	return ingredients, nil
}
