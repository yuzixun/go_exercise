package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func handleDownload(fetchChan chan Fetcher) {
	for {
		select {
		case data := <-fetchChan:
			fmt.Println(data)
			downloadFile(data)
		}
	}
}

func downloadFile(fetchChan Fetcher) error {
	url := fetchChan.url
	filepath := fetchChan.fileName

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
