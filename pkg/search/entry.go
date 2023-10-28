package search

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/duosonic62/aozora-search/pkg/domain"
	"golang.org/x/text/encoding/japanese"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

func FindEntries(siteURL string) ([]domain.Entry, error) {
	response, err := http.Get(siteURL)
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
	entries := []domain.Entry{}
	doc.Find("ol li a").Each(func(i int, selection *goquery.Selection) {
		title := selection.Text()
		token := pat.FindStringSubmatch(selection.AttrOr("href", ""))
		if len(token) != 3 {
			return
		}
		pageURL := fmt.Sprintf("https://www.aozora.gr.jp/cards/%s/card%s.html", token[1], token[2])
		author, zipURL := findAuthorAndZIP(pageURL)
		if zipURL != "" {
			entries = append(entries, domain.Entry{
				AuthorID: token[1],
				Author:   author,
				TitleID:  token[2],
				Title:    title,
				InfoURL:  siteURL,
				ZipURL:   zipURL,
			})
		}
	})

	return entries, nil
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

func ExtractText(zipURL string) (string, error) {
	resp, err := http.Get(zipURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return "", err
	}

	for _, file := range r.File {
		if path.Ext(file.Name) == ".txt" {
			f, err := file.Open()
			if err != nil {
				return "", err
			}
			b, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				return "", err
			}
			b, err = japanese.ShiftJIS.NewDecoder().Bytes(b)
			if err != nil {
				return "", err
			}

			return string(b), nil
		}
	}

	return "", errors.New("contents not found")
}
