package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func CopyFile(inPath string, outPath string) error {
	input, err := ioutil.ReadFile(inPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = ioutil.WriteFile(outPath, input, 0644)
	if err != nil {
		fmt.Println("Error creating", outPath)
		fmt.Println(err)
		return err
	}

	return nil
}

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	return nil
}
