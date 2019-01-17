package main

import (
	"fmt"
	"github.com/disposedtrolley/newsbuddy/pkg/formatter"
	"github.com/disposedtrolley/newsbuddy/pkg/models"
	"github.com/disposedtrolley/newsbuddy/pkg/parser"
	"github.com/disposedtrolley/newsbuddy/pkg/summariser"
	"github.com/disposedtrolley/newsbuddy/pkg/util"
	"github.com/disposedtrolley/newsbuddy/pkg/writer"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("No input file, template, or output directory received. Usage: newsbuddy <input file> <template> <output dir>")
	}

	filePath := os.Args[1]
	templatePath := os.Args[2]
	outDir := os.Args[3]

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

		go func(url string, articleType string, category string, summary string) {
			log.Printf("[main async] Fetching title for article %s\n", url)
			title, err := summariser.ArticleTitle(url)

			if err != nil {
				fmt.Println(err)
			}

			ch <- &models.Article{URL: url, Title: title, Type: articleType, Category: category, Summary: summary}
		}(url, article.Type, article.Category, article.Summary)
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

			outStr, err := formatter.FillTemplate(templatePath, data)

			if err != nil {
				log.Fatalf("[main] An error occurred when generating the template: %s\n", err.Error())
			}

			// Output
			log.Printf("[main] Writing output to file...")

			formattedDate := time.Now().Local().Format("2006-01-02")

			outDir := fmt.Sprintf("%s/%s", outDir, formattedDate)

			util.CreateDirIfNotExists(outDir)

			// Write filled template
			w := writer.NewFileWriter(fmt.Sprintf("%s/filled_template.mjml", outDir))
			w.WriteToFile(outStr)
			// Copy input source file
			util.CopyFile(filePath, fmt.Sprintf("%s/source.toml", outDir))
			// Invoke mjml binary to output an HTML file
			util.RunExternalCommand(
				"./mjml/node_modules/.bin/mjml",
				[]string{fmt.Sprintf("%s/filled_template.mjml", outDir), "-o", fmt.Sprintf("%s/out.html", outDir)})

			return
		}
	}
}
