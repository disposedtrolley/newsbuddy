package formatter

import (
	"bytes"
	"log"
	"text/template"
)

type Article struct {
	URL   string
	Title string
	Type  string
}

type NewsletterData struct {
	Title       string
	IssueNo     int
	PubDate     string
	WelcomeText string
	Articles    []Article
}

func FillTemplate(data NewsletterData) (string, error) {
	t := template.Must(template.ParseFiles("./template.mjml"))

	buf := new(bytes.Buffer)

	err := t.Execute(buf, data)

	if err != nil {
		log.Printf("[formatter] An error occurred when executing the template: %s\n", err.Error())
	}

	return buf.String(), err
}
