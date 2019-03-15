package summariser

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/disposedtrolley/newsbuddy/pkg/models"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// HostImage uploads an image at the given path to Imgur
// and returns the hosted URL.
func HostImage(filePath string) (models.ImgurResponse, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(filePath)
	if err != nil {
		return models.ImgurResponse{}, err
	}
	defer f.Close()
	fw, err := w.CreateFormFile("image", filePath)
	if err != nil {
		return models.ImgurResponse{}, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return models.ImgurResponse{}, err
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", &b)
	if err != nil {
		return models.ImgurResponse{}, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Client-ID 18b425710d7148c")

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return models.ImgurResponse{}, err
	}

	var imgurResponse models.ImgurResponse
	strBody, err := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(strBody), &imgurResponse)

	return imgurResponse, err
}

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
					elementValue := strings.TrimSpace(z.Token().Data)
					return elementValue, nil
				}
			}
		}
	}
}
