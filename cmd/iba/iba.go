package iba

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"net/http"

	"github.com/ramonmedeiros/iba/cmd/model"
)

const cocktailSitemap = "https://iba-world.com/wp-sitemap-posts-iba-cocktail-1.xml"

type ibaParser struct {
}

func New() *ibaParser {
	return &ibaParser{}
}

func (p *ibaParser) GetSource() string {
	return "iba"
}

func (p *ibaParser) GetLinks() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, cocktailSitemap, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var newSitemap model.Sitemap
	err = xml.NewDecoder(bytes.NewBufferString(string(content))).Decode(&newSitemap)
	if err != nil {
		log.Fatal(err)
	}

	var urls []string
	for _, url := range newSitemap.URL {
		urls = append(urls, url.Loc)
	}

	return urls, nil
}
