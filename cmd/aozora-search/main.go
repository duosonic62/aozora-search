package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var dsn string
	flag.StringVar(&dsn, "d", "database.sqlite", "database")
	flag.Usage = func() {
		fmt.Println("usage")
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "authors":
		//err = showAuthors(db)
		println("show authors")
	case "titles":
		if flag.NArg() != 2 {
			flag.Usage()
			os.Exit(2)
		}
		//err = showTitles(db, flagArg(1))
		println("show titles")
	case "content":
		if flag.NArg() != 3 {
			flag.Usage()
			os.Exit(2)
		}
		//err = showContent(db, flagArg(1), flagArg(2))
		println("show content")
	case "query":
		if flag.NArg() != 2 {
			flag.Usage()
			os.Exit(2)
		}
		//err = queryContent(db, flagArg(1))
		println("query Content")
	}

	if err != nil {
		log.Fatal(err)
	}
}
