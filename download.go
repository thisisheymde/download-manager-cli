package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type DownloadManager struct {
	URL       string
	OutputDir string
	NumParts  int64
	Client    *http.Client
	Parts     [][2]int64
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

	tr := &http.Transport{
		ResponseHeaderTimeout: 15 * time.Second,
	}

	return &DownloadManager{
		URL:       url,
		OutputDir: outputdir,
		NumParts:  numParts,
		Client: &http.Client{
			Transport: tr,
		},
		Parts: segments,
	}, nil
}

func (dm *DownloadManager) Download(count int) error {

	timeTaken := time.Now()

	out, err := CreateFile(dm.URL, dm.OutputDir)
	if err != nil {
		return err
	}
	defer out.Close()

	var wg sync.WaitGroup

	errChan := make(chan error)

	for i, segment := range dm.Parts {
		wg.Add(1)
		go func(i int, segment [2]int64) {
			defer wg.Done()

			req, err := http.NewRequest("GET", dm.URL, nil)
			if err != nil {
				errChan <- err
				return
			}

			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", segment[0], segment[1]))

			resp, err := dm.Client.Do(req)
			if err != nil {
				errChan <- err
				return
			}

			defer resp.Body.Close()

			// download

			buf := make([]byte, 1024)

			for {
				n, err := resp.Body.Read(buf)
				if err != nil && err != io.EOF {
					errChan <- err
					return
				}
				if n == 0 {
					break
				}
				_, err = out.WriteAt(buf[:n], segment[0])
				if err != nil {
					errChan <- err
					return
				}
				segment[0] += int64(n)
			}

		}(i, segment)
	}

	select {
	case err := <-errChan:
		DeleteFile(dm.URL, dm.OutputDir)
		return err

	default:
		wg.Wait()

		fmt.Printf("Download %v Completed! took : %v \n", count, time.Since(timeTaken))

		return nil
	}

}
