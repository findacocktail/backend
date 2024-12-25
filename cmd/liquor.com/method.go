package liquorcom

import (
	"github.com/findacocktail/backend/cmd/parsing"
	"golang.org/x/net/html"
)

func parseMethod(token *html.Node) (string, error) {
	instructions, err := parsing.GetNodeByAttribute(token, "id", "section--instructions_1-0")
	if err != nil {
		return "", err
	}

	instructionsList, err := parsing.GetNode(instructions, "ol")
	if err != nil {
		return "", err
	}

	method := ""
	for node := range instructionsList.ChildNodes() {
		if node.Data != "li" {
			continue
		}

		step, err := parsing.GetNode(node, "p")
		if err != nil {
			return "", err
		}

		method += step.FirstChild.Data + "\n"
	}

	return method, nil
}
