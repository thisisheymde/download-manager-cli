package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type DownloadManager struct {
	URL       string
	OutputDir string
	NumParts  int64
	Client    *http.Client
	Parts     [][2]int64
	Progress  []int64
}

func NewDownloadManager(url string, outputdir string, numParts int64) (*DownloadManager, error) {
	resp, err := http.Head(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	segmentSize := resp.ContentLength / numParts
	segments := make([][2]int64, numParts)

	for i := range segments {
		start := int64(i) * segmentSize
		end := start + segmentSize - 1
		if i == int(numParts-1) {
			end = resp.ContentLength - 1
		}
		segments[i] = [2]int64{start, end}
	}

	return &DownloadManager{
		URL:       url,
		OutputDir: outputdir,
		NumParts:  numParts,
		Client:    &http.Client{},
		Parts:     segments,
		Progress:  make([]int64, numParts),
	}, nil
}

// fix error handling
func (dm *DownloadManager) Download() error {

	// create file
	out, err := CreateFile(dm.URL, dm.OutputDir)
	if err != nil {
		return err
	}
	defer out.Close()

	var wg sync.WaitGroup

	for i, segment := range dm.Parts {
		wg.Add(1)
		go func(i int, segment [2]int64) {
			defer wg.Done()

			req, err := http.NewRequest("GET", dm.URL, nil)
			if err != nil {
				log.Println(err)
				return
			}

			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", segment[0], segment[1]))

			resp, err := dm.Client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}

			defer resp.Body.Close()

			// download

			buf := make([]byte, 1024)

			for {
				n, err := resp.Body.Read(buf)
				if err != nil && err != io.EOF {
					log.Println(err)
					return
				}
				if n == 0 {
					break
				}
				_, err = out.WriteAt(buf[:n], segment[0])
				if err != nil {
					log.Println(err)
					return
				}
				segment[0] += int64(n)

				dm.Progress[i] = dm.Progress[i] + 1
			}

		}(i, segment)
	}

	wg.Wait()

	log.Println("Download Completed!")

	return nil
}
