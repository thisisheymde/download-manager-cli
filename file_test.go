package main

import (
	"os"
	"testing"
)

func TestCreateFile(t *testing.T) {
	filedir := "./"
	filename := "test.txt"
	filePath := filedir + filename

	// Create a new file
	file, err := CreateFile("https://example.com/test.txt", filedir)
	if err != nil {
		t.Errorf("CreateFile Error: %s", err.Error())
	}
	defer file.Close()

	// Verify that the file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("CreateFile did not create the file")
	}

	// Try to create a file that already exists
	_, err = CreateFile("https://example.com/test.txt", filedir)
	if err == nil {
		t.Error("No error for creating an already existing file")
	}

	// Clean up file
	os.Remove(filePath)
}

func TestDeleteFile(t *testing.T) {
	filedir := "./"
	filename := "test.txt"
	filePath := filedir + filename

	// Create a file to be deleted
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file for test: %s", err.Error())
	}
	file.Close()

	// Delete an existing file
	DeleteFile("https://example.com/test.txt", filedir)

	// Verify that the file was deleted
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("DeleteFile did not delete the file")
	}
}
