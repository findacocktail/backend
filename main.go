package main

import "fmt"

func main() {
	urls, err := parseSitemap()
	if err != nil {
		panic(err)
	}

	fmt.Println(urls)
}
