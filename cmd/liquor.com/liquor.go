package liquorcom

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"net/http"

	"github.com/findacocktail/backend/cmd/model"
)

const cocktailSitemap = "https://www.liquor.com/sitemap_1.xml"

type liquorParser struct {
}

func New() *liquorParser {
	return &liquorParser{}
}

func (p *liquorParser) GetSource() string {
	return "liquor"
}

func (p *liquorParser) GetLinks() ([]string, error) {
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
