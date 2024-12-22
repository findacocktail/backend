package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Recipe struct {
	Name        string        `json:"name"`
	YoutubeLink string        `json:"youtube_link"`
	Ingredients []*Ingredient `json:"ingredients"`
	Method      string        `json:"method"`
	Garnish     string        `json:"garnish"`
	ImageURL    string        `json:"image_url"`
}

type Ingredient struct {
	Amount      float64 `json:"amount"`
	Scale       string  `json:"scale"`
	Description string  `json:"description"`
}

func parseRecipe(recipeLink string) (*Recipe, error) {
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

	cocktailName, err := getNode(token, "h1")
	if err != nil {
		return nil, err
	}

	youtubeLink, err := getLinkStartsWith(token, "https://www.youtube.com/watch")
	if err != nil {
		fmt.Println("could not find link", err)
	}

	recipe := Recipe{
		Name:        strings.TrimSpace(cocktailName.FirstChild.Data),
		YoutubeLink: youtubeLink,
	}

	image, err := getNodeByAttribute(token, "fetchpriority", "high")
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

func getNode(root *html.Node, tagText string) (*html.Node, error) {
	var tag *html.Node
	var f func(node *html.Node)
	f = func(node *html.Node) {
		if strings.EqualFold(node.Data, tagText) {
			tag = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(root)
	if tag != nil {
		return tag, nil
	}
	return nil, errors.New(fmt.Sprint("tag not found", tagText))
}

func getNodeByAttribute(root *html.Node, attributeName string, attributeValue string) (*html.Node, error) {
	var tag *html.Node
	var f func(node *html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode {

			for _, attr := range node.Attr {
				if attr.Key == attributeName &&
					attr.Val == attributeValue {
					tag = node
					return
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(root)
	if tag != nil {
		return tag, nil
	}
	return nil, errors.New(fmt.Sprint("tag not found", attributeName, attributeName))
}

func getLinkStartsWith(root *html.Node, link string) (string, error) {
	var f func(node *html.Node)
	var href string
	f = func(node *html.Node) {
		if node.Type == html.ElementNode {

			for _, attr := range node.Attr {
				if attr.Key == "href" &&
					strings.HasPrefix(attr.Val, link) {
					href = attr.Val
					return
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(root)
	if href != "" {
		return href, nil
	}
	return href, errors.New(fmt.Sprint("tag not found", link))
}
