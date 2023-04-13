package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sync"
)

func main() {
	var url string
	var dir string
	var connex int64
	var urllist string

	flag.StringVar(&url, "url", "", "put url in here")
	flag.StringVar(&dir, "dir", os.Getenv("HOME")+"/Downloads/", "put dir in here")
	flag.Int64Var(&connex, "connections", 4, "put no of connections in here")
	flag.StringVar(&urllist, "multiple", "", "put filename containing the urls here")

	// Parse the command-line arguments
	flag.Parse()

	if url != "" && urllist == "" {
		downloader, err := NewDownloadManager(url, dir, connex)

		if err != nil {
			panic(err)
		}

		downloader.Download()
	} else if urllist != "" && url == "" {
		file, err := os.Open(urllist)
		if err != nil {
			log.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		var mwg sync.WaitGroup

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			mwg.Add(1)
			go func(d string) {
				defer mwg.Done()

				downloader, err := NewDownloadManager(d, dir, 6)

				if err != nil {
					panic(err)
				}

				downloader.Download()
			}(line)
		}

		if err := scanner.Err(); err != nil {
			log.Println("Error reading file:", err)
		}

		mwg.Wait()

	} else {
		log.Println("enter flags correctly.")
	}

	log.Println("Finally Free!!!")
}
