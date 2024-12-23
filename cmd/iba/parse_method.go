package iba

import (
	"strings"

	"golang.org/x/net/html"
)

func parseListOfP(root *html.Node, header string) (string, error) {
	node, err := getNode(root, header)
	if err != nil {
		return "", err
	}

	topParent := node.Parent.Parent.Parent
	nextNode := topParent.NextSibling.NextSibling

	div, err := getNodeByAttribute(nextNode, "class", "elementor-shortcode")
	if err != nil {
		return "", err
	}

	p := div.FirstChild

	method := []string{}
	for p != nil {
		if p.FirstChild != nil &&
			p.FirstChild.Type == html.TextNode {
			method = append(method, p.FirstChild.Data)
		}
		p = p.NextSibling
	}

	return strings.Join(method, " "), nil
}
