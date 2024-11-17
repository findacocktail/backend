package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

const cocktailSitemap = "https://iba-world.com/wp-sitemap-posts-iba-cocktail-1.xml"

type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URL     []struct {
		Loc     string `xml:"loc"`
		Lastmod string `xml:"lastmod"`
	} `xml:"url"`
}

func parseSitemap() ([]string, error) {
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

	fmt.Println(string(content))

	var newSitemap Urlset
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
