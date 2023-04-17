package main

import (
	"bufio"
	"flag"
	"fmt"
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

	//Parse the command-line arguments
	flag.Parse()

	count := 0

	if url != "" && urllist == "" {
		downloader, err := NewDownloadManager(url, dir, connex)

		if err != nil {
			panic(err)
		}

		err = downloader.Download(count)

		if err != nil {
			panic(err)
		}

	} else if urllist != "" && url == "" {
		file, err := os.Open(urllist)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		var mwg sync.WaitGroup

		scanner := bufio.NewScanner(file)
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}

		for scanner.Scan() {
			line := scanner.Text()
			mwg.Add(1)
			go func(d string) {
				defer mwg.Done()

				downloader, err := NewDownloadManager(d, dir, 6)

				if err != nil {
					panic(err)
				}

				count += 1
				err = downloader.Download(count)

				if err != nil {
					panic(err)
				}

			}(line)
		}

		mwg.Wait()

	} else {
		fmt.Println("enter flags correctly.")
	}
}
