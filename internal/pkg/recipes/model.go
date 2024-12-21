package recipes

type Recipe struct {
	Name        string        `json:"name"`
	YoutubeLink string        `json:"youtube_link"`
	Ingredients []*Ingredient `json:"ingredients"`
	Method      string        `json:"method"`
	Garnish     string        `json:"garnish"`
}

type Ingredient struct {
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Scale       string  `json:"scale"`
	Description string  `json:"description"`
}
