package formatter

import (
	"bytes"
	"github.com/disposedtrolley/newsbuddy/pkg/models"
	"log"
	"text/template"
)

func FillTemplate(templatePath string, data models.NewsletterData) (string, error) {
	t := template.Must(template.ParseFiles(templatePath))

	buf := new(bytes.Buffer)

	err := t.Execute(buf, data)

	if err != nil {
		log.Printf("[formatter] An error occurred when executing the template: %s\n", err.Error())
	}

	return buf.String(), err
}
