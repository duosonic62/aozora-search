package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

type Entry struct {
	AuthorID string
	Author   string
	TitleID  string
	Title    string
	InfoURL  string
	ZipURL   string
}

func findEntries(sitURL string) ([]Entry, error) {
	response, err := http.Get(sitURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
	doc.Find("ol li a").Each(func(i int, selection *goquery.Selection) {
		token := pat.FindStringSubmatch(selection.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}
		pageURL := fmt.Sprintf("https://www.aozora.gr.jp/cards/%s/card%s.html", token[1], token[2])
		author, zipURL := findAuthorAndZIP(pageURL)
		println(author, zipURL)
		println(pageURL)
	})

	return nil, nil
}

func findAuthorAndZIP(siteURL string) (string, string) {
	response, err := http.Get(siteURL)
	if err != nil {
		log.Fatal(err)
		return "", ""
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
		return "", ""
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)

	author := doc.Find("table[summary=作家データ] tr:nth-child(2) td:nth-child(2)").Text()

	zipURL := ""
	doc.Find("table.download a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")
		if strings.HasSuffix(href, ".zip") {
			zipURL = href
		}
	})
	if zipURL == "" {
		return author, ""
	}
	if strings.HasPrefix(zipURL, "https://") || strings.HasPrefix(zipURL, "http://") {
		return author, zipURL
	}
	u, err := url.Parse(siteURL)
	if err != nil {
		return author, ""
	}
	u.Path = path.Join(path.Dir(u.Path), zipURL)

	return author, u.String()
}

func main() {
	listURL := "https://www.aozora.gr.jp/index_pages/person879.html"

	entries, err := findEntries(listURL)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Println(entry.Title, entry.ZipURL)
	}
}
