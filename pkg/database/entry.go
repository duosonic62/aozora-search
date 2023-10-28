package database

import (
	"database/sql"
	"github.com/duosonic62/aozora-search/pkg/domain"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"strings"
)

func AddEntry(db *sql.DB, entry *domain.Entry, content string) error {
	_, err := db.Exec(`REPLACE INTO authors(author_id, author) values(?, ?)`, entry.AuthorID, entry.Author)
	if err != nil {
		return err
	}

	res, err := db.Exec(`REPLACE INTO contents(author_id, title_id, title, content) values (?, ?, ?, ?)`,
		entry.AuthorID,
		entry.TitleID,
		entry.Title,
		content,
	)
	if err != nil {
		return err
	}
	docID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return err
	}

	seq := t.Wakati(content)
	_, err = db.Exec(`REPLACE INTO contents_fts(docid, words) values(?, ?)`, docID, strings.Join(seq, " "))
	if err != nil {
		return err
	}

	return nil
}
