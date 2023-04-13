package main

import (
	"os"
	"regexp"
)

const pattern = `[^\/]+$`

func CreateFile(url string, filedir string) (*os.File, error) {

	re := regexp.MustCompile(pattern)
	filename := re.FindString(url)

	out, err := os.Create(filedir + filename)
	if err != nil {
		return nil, err
	}

	return out, nil
}
