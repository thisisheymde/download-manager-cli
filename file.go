package main

import (
	"errors"
	"os"
	"regexp"
)

const pattern = `[^\/]+$`

func CreateFile(url string, filedir string) (*os.File, error) {

	re := regexp.MustCompile(pattern)
	filename := re.FindString(url)

	_, err := os.Stat(filename)

	if !os.IsNotExist(err) {

		return nil, errors.New("file already exists")

	} else {
		out, err := os.Create(filedir + filename)
		if err != nil {
			return nil, errors.New("failed to create file")
		}

		return out, nil
	}

}
