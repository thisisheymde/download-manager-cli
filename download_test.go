package main

import (
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	url := "http://speedtest.ftp.otenet.gr/files/test10Mb.db"
	dir := "./"
	filename := "test10Mb.db"
	filePath := dir + filename

	downloader, err := NewDownloadManager(url, dir, 6)

	if err != nil {
		t.Errorf("NewDownloadManager Error: %s", err.Error())
	}

	err = downloader.Download(1)

	if err != nil {
		t.Errorf("Download Error: %s", err.Error())
	}

	// Verify if the file has been downloaded
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Download did not download the file")
	}

	// Clean up file
	os.Remove(filePath)
}
