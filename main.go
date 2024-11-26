package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	concurrentRequest = 4
)

func main() {
	urls, err := parseSitemap()
	if err != nil {
		panic(err)
	}

	start := 0
	run := test{
		recipes: []*Recipe{},
		wg:      &sync.WaitGroup{},
	}
	totalURLs := len(urls)
	for start < totalURLs {
		diff := totalURLs - start
		if concurrentRequest > diff {
			concurrentRequest = diff
		}

		for range concurrentRequest {
			run.dispatchJob(urls[start])
			start++
		}
		run.wg.Wait()
	}

	content, err := json.Marshal(run.recipes)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}

type test struct {
	wg      *sync.WaitGroup
	recipes []*Recipe
}

// Search for prime numbers in 4 ranges.
func (t *test) dispatchJob(url string) {
	t.wg.Add(1)
	go func(url string) {
		recipe, err := parseRecipe(url)
		if err != nil {
			panic(err)
		}
		t.recipes = append(t.recipes, recipe)

		t.wg.Done()
	}(url)
}
