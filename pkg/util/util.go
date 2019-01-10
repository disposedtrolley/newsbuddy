package util

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func CopyFile(inPath string, outPath string) error {
	input, err := ioutil.ReadFile(inPath)
	if err != nil {
		log.Printf("Error copying file from %s to %s: %s\n", inPath, outPath, err)
		return err
	}

	err = ioutil.WriteFile(outPath, input, 0644)
	if err != nil {
		log.Printf("Error creating directory %s: %s\n", outPath, err)
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

func RunExternalCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	log.Printf("Running command and waiting for it to finish...")

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return err
}
