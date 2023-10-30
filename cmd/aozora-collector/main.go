package main

import (
	"github.com/duosonic62/aozora-search/pkg/database"
	"github.com/duosonic62/aozora-search/pkg/search"
	"log"
)

func main() {
	db, err := database.SetupDB("database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := search.FindEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found %d entries", len(entries))
	for _, entry := range entries {
		log.Printf("adding %+v\n", entry)
		content, err := search.ExtractText(entry.ZipURL)
		if err != nil {
			log.Println(err)
			continue
		}
		err = database.AddEntry(db, &entry, content)
		if err != nil {
			log.Println(err)
		}
	}
}
