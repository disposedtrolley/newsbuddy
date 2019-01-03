package parser

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// ReadFileIntoArray reads a .txt file, and returns an array of
// strings representing each non-empty line of the file.
func ReadFileIntoArray(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	var lines []string

	for _, line := range strings.Split(string(content), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, err
}
