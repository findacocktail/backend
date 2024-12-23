package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ramonmedeiros/iba/cmd/iba"
	"github.com/ramonmedeiros/iba/cmd/model"
)

var (
	concurrentRequest = 4
)

func Scrape() {
	for _, parser := range []model.Parser{iba.New()} {
		urls, err := parser.GetLinks()
		if err != nil {
			panic(err)
		}

		start := 0
		run := test{
			recipes: []*model.Recipe{},
			wg:      &sync.WaitGroup{},
			parser:  parser,
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

		now := time.Now()
		fileName := fmt.Sprintf("%d%d%d_%d%d_%s.json",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			parser.GetSource(),
		)

		err = os.WriteFile(fileName, content, os.ModeAppend)
		if err != nil {
			panic(err)
		}
	}

}

type test struct {
	wg      *sync.WaitGroup
	parser  model.Parser
	recipes []*model.Recipe
}

// Search for prime numbers in 4 ranges.
func (t *test) dispatchJob(url string) {
	t.wg.Add(1)
	go func(url string) {
		recipe, err := t.parser.GetRecipe(url)
		if err != nil {
			panic(err)
		}
		t.recipes = append(t.recipes, recipe)

		t.wg.Done()
	}(url)
}
