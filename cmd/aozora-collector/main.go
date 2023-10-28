package main

import (
	"fmt"
	"github.com/duosonic62/aozora-search/pkg/search"
	"log"
)

func main() {
	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := search.FindEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		content, err := search.ExtractText(entry.ZipURL)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(entry.Title, entry.ZipURL)
		fmt.Println(content)
	}
}
