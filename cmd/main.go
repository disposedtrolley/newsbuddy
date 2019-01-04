package main

import (
	"fmt"
	"github.com/disposedtrolley/newsbuddy/pkg/formatter"
	"github.com/disposedtrolley/newsbuddy/pkg/models"
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

	source, err := parser.ReadIssueSource(filePath)

	if err != nil {
		fmt.Printf("[error] Unable to read input file: %s\n", err.Error())
	}

	// Main processing loop
	var articles []models.Article
	ch := make(chan *models.Article, len(source.Articles))

	for i, article := range source.Articles {
		url := article.Url
		log.Printf("[main] Processing article %d of %d\n", i+1, len(source.Articles))

		go func(url string, articleType string, category string) {
			log.Printf("[main async] Fetching title for article %s\n", url)
			title, err := summariser.ArticleTitle(url)

			if err != nil {
				fmt.Println(err)
			}

			ch <- &models.Article{URL: url, Title: title, Type: articleType, Category: category}
		}(url, article.Type, article.Category)
	}

	// Iterate through the buffered channel and append processed articles
	// to the articles slice.
	for r := range ch {
		articles = append(articles, *r)

		// Once all articles have been processed, it's time to write to the
		// template.
		if len(articles) == len(source.Articles) {
			log.Printf("[main] Articles processed. Assembling data required for the template...")

			data := models.NewsletterData{
				Title:       source.Metadata.Title,
				IssueNo:     source.Metadata.IssueNo,
				PubDate:     source.Metadata.PubDate,
				WelcomeText: source.Metadata.WelcomeText,
				Articles:    articles}

			outStr, err := formatter.FillTemplate(data)

			if err != nil {
				log.Fatalf("[main] An error occurred when generating the template: %s\n", err.Error())
			}

			log.Printf("[main] Writing output to file...")

			outFileName := fmt.Sprintf("%s.mjml", source.Metadata.PubDate)

			w := writer.NewFileWriter(fmt.Sprintf(outFileName, filePath))

			w.WriteToFile(outStr)

			return
		}
	}
}
