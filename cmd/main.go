package main

import (
	"fmt"
	"github.com/disposedtrolley/newsbuddy/pkg/formatter"
	"github.com/disposedtrolley/newsbuddy/pkg/parser"
	"github.com/disposedtrolley/newsbuddy/pkg/summariser"
	"github.com/disposedtrolley/newsbuddy/pkg/writer"
	"log"
	"os"
)

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
	var articles []formatter.Article
	for i, url := range urls {
		log.Printf("[main] Processing article %d of %d\n", i+1, len(urls))
		title, err := summariser.ArticleTitle(url)

		if err != nil {
			log.Printf("[main] Error encountered while processing %s: %s\n", url, err.Error())
		} else {
			articles = append(articles, formatter.Article{URL: url, Title: title, Type: "TEXT"})
		}
	}

	log.Printf("[main] Articles processed. Assembling data required for the template...")

	data := formatter.NewsletterData{
		Title:       "iX Friday Fill-Ins",
		IssueNo:     15,
		PubDate:     "04 Jan 2019",
		WelcomeText: "Hello",
		Articles:    articles}

	outStr, err := formatter.FillTemplate(data)

	if err != nil {
		log.Fatalf("[main] An error occurred when generating the template: %s\n", err.Error())
	}

	log.Printf("[main] Writing output to file...")

	w := writer.NewFileWriter(fmt.Sprintf("%s.mjml", filePath))

	w.WriteToFile(outStr)
}
