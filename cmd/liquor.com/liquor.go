package liquorcom

import (
	"bytes"
	"embed"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/findacocktail/backend/cmd/model"
)

//go:embed cache/sitemap_1.xml
var localFile embed.FS

const cocktailSitemap = "https://www.liquor.com/sitemap_1.xml"

type liquorParser struct {
	cache bool
}

func New(cache bool) *liquorParser {
	return &liquorParser{
		cache: cache,
	}
}

func (p *liquorParser) GetSource() string {
	return "liquor"
}

func (p *liquorParser) GetLinks() ([]string, error) {
	var content []byte
	var err error
	if p.cache {
		content, err = localFile.ReadFile("cache/sitemap_1.xml")
		if err != nil {
			return nil, err
		}
	} else {
		req, err := http.NewRequest(http.MethodGet, cocktailSitemap, nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		content, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	var newSitemap model.Sitemap
	err = xml.NewDecoder(bytes.NewBufferString(string(content))).Decode(&newSitemap)
	if err != nil {
		log.Fatal(err)
	}

	var urls []string
	for _, url := range newSitemap.URL {
		if strings.HasPrefix(url.Loc, "https://www.liquor.com/recipes/") {
			urls = append(urls, url.Loc)
		}
	}

	return urls, nil
}
