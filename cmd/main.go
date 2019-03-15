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

	// Track the number of async tasks to process.
	asyncTasks := len(source.Articles) + len(source.Images)

	// Main processing loop
	var articles []models.Article
	var images []models.Image
	ch := make(chan interface{}, asyncTasks)

	// Asynchronously process articles
	for i, article := range source.Articles {
		log.Printf("[main] Processing article %d of %d\n", i+1, len(source.Articles))

		go func(article models.Article) {
			log.Printf("[main async] Fetching title for article %s\n", article.URL)
			title, err := summariser.ArticleTitle(article.URL)

			if err != nil {
				fmt.Println(err)
			}

			article.Title = title

			ch <- &article
		}(article)
	}

	// Asynchronously process images
	for i, image := range source.Images {
		log.Printf("[main] Processing image %d of %d\n", i+1, len(source.Images))

		go func(image models.Image) {
			log.Printf("[main async] Uploading image at path %s\n", image.FilePath)
			uploadedImage, err := summariser.HostImage(image.FilePath)

			if err != nil {
				fmt.Println(err)
			}

			image.URL = uploadedImage.Data.Link
			image.Width = uploadedImage.Data.Width

			ch <- &image
		}(image)
	}

	// Iterate through the buffered channel and append processed articles
	// and images to their respective slices.
	for processedElement := range ch {
		switch processedElement := processedElement.(type) {
		case *models.Article:
			articles = append(articles, *processedElement)
			break
		case *models.Image:
			images = append(images, *processedElement)
			break
		default:
			log.Printf("[main async] Processed task is the wrong type.")
		}

		// Once all async tasks have been processed, it's time to write to the
		// template.
		if len(articles)+len(images) == asyncTasks {
			log.Printf("[main] Tasks processed. Assembling data required for the template...")

			data := models.NewsletterData{
				Title:       source.Metadata.Title,
				IssueNo:     source.Metadata.IssueNo,
				PubDate:     source.Metadata.PubDate,
				WelcomeText: source.Metadata.WelcomeText,
				Articles:    articles,
				Images:      images}

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
