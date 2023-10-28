package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func setupDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	queries := []string{
		`CREATE TABLE IF NOT EXISTS author(author_id TEXT, author TEXT, PRIMARY KEY (author_id))`,
		`CREATE TABLE IF NOT EXISTS contents(author_id TEXT, title_id TEXT, title TEXT, PRIMARY KEY (author_id, title_id))`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS contents_fts USING fts4(words)`,
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	return db, nil
}
