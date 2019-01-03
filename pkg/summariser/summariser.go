package summariser

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// ArticleTitle extracts the title of the HTML page at a given URL.
func ArticleTitle(url string) (string, error) {
	log.Printf("[summariser] Extracting article title for %s ...\n", url)
	resp, err := http.Get(url)

	// handle the error if there is one
	if err != nil {
		log.Printf("[summariser] Unable to GET article: %s\n", err.Error())

		return "", err
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[summariser] Unable to read article body: %s\n", err.Error())
	}

	body := strings.NewReader(string(html))

	title, err := extractElementValue(body, "title")

	if err != nil {
		log.Printf("[summariser] Unable to extract title from article: %s\n", err.Error())
	}

	return title, err
}

func extractElementValue(body io.Reader, tag string) (string, error) {
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return "", errors.New("Unable to parse HTML")
		} else if tt == html.StartTagToken {
			el := z.Token()
			if el.Data == tag {
				if tt = z.Next(); tt == html.TextToken {
					articleTitle := strings.TrimSpace(z.Token().Data)
					return articleTitle, nil
				}
			}
		}
	}
}
