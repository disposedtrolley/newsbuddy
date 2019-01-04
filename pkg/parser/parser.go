package parser

import (
	"github.com/BurntSushi/toml"
	"github.com/disposedtrolley/newsbuddy/pkg/models"
	"io/ioutil"
	"log"
	"strings"
)

func readFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("[parser] Unable to read file at path %s: %s\n", path, err.Error())
	}

	return content, err
}

// ReadFileIntoArray reads a .txt file, and returns an array of
// strings representing each non-empty line of the file.
func ReadFileIntoArray(path string) ([]string, error) {
	content, err := readFile(path)

	var lines []string

	for _, line := range strings.Split(string(content), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, err
}

func ReadIssueSource(path string) (models.SourceFile, error) {
	content, err := readFile(path)

	if err != nil {
		log.Fatal(err)
	}

	var data models.SourceFile
	if _, err := toml.Decode(string(content), &data); err != nil {
		log.Printf("[parser] Unable to decode TOML file: %s\n", err.Error())
	}

	return data, err
}
