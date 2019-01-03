package writer

import (
	"fmt"
	"log"
	"os"
)

// FileWriter provides utilities to write strings to a file.
type FileWriter struct {
	Filename string
}

// NewFileWriter instantiates a new FileWriter, setting the
// Filename property.
func NewFileWriter(outFile string) *FileWriter {
	return &FileWriter{Filename: outFile}
}

// WriteToFile writes a string to the file.
func (fw FileWriter) WriteToFile(text string) {
	file, err := os.Create(fw.Filename)
	if err != nil {
		log.Fatal("[writer] Unable to create destination file: %s\n", err.Error())
	}
	defer file.Close()

	fmt.Fprint(file, text)
}
