package parsing

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func GetNode(root *html.Node, tagText string) (*html.Node, error) {
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
	return nil, errors.New(fmt.Sprint("tag not found ", tagText))
}

func GetNodeByAttribute(root *html.Node, attributeName string, attributeValue string) (*html.Node, error) {
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
	return nil, fmt.Errorf("tag not found %s %s", attributeName, attributeValue)
}

func GetAttributeStartsWith(root *html.Node, attributeName string, link string) (string, error) {
	var f func(node *html.Node)
	var attribute string
	f = func(node *html.Node) {
		if node.Type == html.ElementNode {

			for _, attr := range node.Attr {
				if attr.Key == attributeName &&
					strings.HasPrefix(attr.Val, link) {
					attribute = attr.Val
					return
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(root)
	if attribute != "" {
		return attribute, nil
	}
	return attribute, errors.New(fmt.Sprint("tag not found ", link))
}
