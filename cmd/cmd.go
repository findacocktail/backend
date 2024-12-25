package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	liquorcom "github.com/findacocktail/backend/cmd/liquor.com"
	"github.com/findacocktail/backend/cmd/model"
)

var (
	concurrentRequest = 4
)

func Scrape() {
	// iba.New()
	for _, parser := range []model.Parser{liquorcom.New(true)} {
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

		err = os.WriteFile(fileName, content, os.ModePerm)
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

func (t *test) dispatchJob(url string) {
	t.wg.Add(1)
	go func(url string) {
		recipe, err := t.parser.GetRecipe(url)
		if err != nil {
			slog.Error("could not fetch", url, err)

		} else {
			t.recipes = append(t.recipes, recipe)
		}

		t.wg.Done()
	}(url)
}
