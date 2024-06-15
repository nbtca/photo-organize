package utils

import (
	"io"
	"log"
	"os"
)

// CopyFile copies a file from src to dst. If src and dst files exist, and are the same, then return success.
// Otherwise, attempt to create a new file, and copy the contents of src to it.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

// write file
func WriteToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return file.Sync()
}
