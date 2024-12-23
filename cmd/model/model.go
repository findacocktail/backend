package model

import "encoding/xml"

type Parser interface {
	GetSource() string
	GetLinks() ([]string, error)
	GetRecipe(recipeLink string) (*Recipe, error)
}

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

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URL     []struct {
		Loc     string `xml:"loc"`
		Lastmod string `xml:"lastmod"`
	} `xml:"url"`
}
