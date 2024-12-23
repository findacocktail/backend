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
	return nil, errors.New(fmt.Sprint("tag not found", tagText))
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
	return nil, errors.New(fmt.Sprint("tag not found", attributeName, attributeName))
}

func GetNodeIfAttributeExists(root *html.Node, attributeName string) (*html.Node, error) {
	var tag *html.Node
	var f func(node *html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode {

			for _, attr := range node.Attr {
				if attr.Key == attributeName {
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

func GetLinkStartsWith(root *html.Node, link string) (string, error) {
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
