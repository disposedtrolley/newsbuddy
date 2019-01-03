package parser

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// ReadFileIntoArray reads a .txt file, and returns an array of
// strings representing each line of the file.
func ReadFileIntoArray(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(content), "\n")
	return lines, err
}
