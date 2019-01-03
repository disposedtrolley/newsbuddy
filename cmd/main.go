package main

import (
	"fmt"
	"github.com/disposedtrolley/newsbuddy/pkg/parser"
	"github.com/disposedtrolley/newsbuddy/pkg/summariser"
	"log"
	"os"
)

type Article struct {
	URL   string
	Title string
	Type  string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No input file received. Usage: newsbuddy <input file>")
	}

	filePath := os.Args[1]

	urls, err := parser.ReadFileIntoArray(filePath)

	if err != nil {
		fmt.Printf("[error] Unable to read input file: %s\n", err.Error())
	}

	// Main processing loop
	var articles []Article
	for _, url := range urls {
		title, err := summariser.ArticleTitle(url)

		if err != nil {
			log.Printf("[main] Error encountered while processing %s: %s\n", url, err.Error())
		} else {
			articles = append(articles, Article{URL: url, Title: title, Type: "TEXT"})
		}
	}

}
