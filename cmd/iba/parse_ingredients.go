package iba

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ramonmedeiros/iba/cmd/model"
	"golang.org/x/net/html"
)

func parseIngredients(root *html.Node) ([]*model.Ingredient, error) {
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
	ingredients := []*model.Ingredient{}
	for child != nil {
		if child.Data == "li" {
			ingredientSplit := strings.Split(child.FirstChild.Data, " ")
			amount, err := strconv.ParseFloat(ingredientSplit[0], 64)
			if err != nil {
				fmt.Println(err)
				amount = 0
				ingredientSplit = []string{"", "", child.FirstChild.Data}
			}

			ingredients = append(ingredients, &model.Ingredient{
				Amount:      amount,
				Scale:       ingredientSplit[1],
				Description: strings.Join(ingredientSplit[2:], " "),
			})
		}
		child = child.NextSibling
	}

	return ingredients, nil
}
