package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func main() {
	urls, err := parseSitemap()
	if err != nil {
		panic(err)
	}

	var recipes []*Recipe
	for _, url := range urls {
		recipe, err := parseRecipe(url)
		if err != nil {
			panic(err)
		}
		recipes = append(recipes, recipe)
	}

	b := []byte("")
	err = json.NewEncoder(bytes.NewBuffer(b)).Encode(recipes)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
