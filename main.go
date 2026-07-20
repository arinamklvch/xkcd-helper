package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"safe_title"`
}

type Response struct {
	str string
	err error
}

func downloadComic(comicNum <-chan int, result chan<- Response) {
	for n := range comicNum {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", n)
		response, err := http.Get(url)
		if err != nil {
			result <- Response{err: fmt.Errorf("Request error: %w", err)}
			return
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusNotFound && n != 404 {
			return
		}
		if response.StatusCode != http.StatusOK {
			result <- Response{err: fmt.Errorf("Unexpected status code: %s", response.Status)}
			return
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			result <- Response{err: fmt.Errorf("Error reading response: %w", err)}
			return
		}
		var comic Comic
		err = json.Unmarshal(body, &comic) // err != nil, err = "Failed unmarshal: unexpected ..."
		if err != nil {
			result <- Response{err: fmt.Errorf("Error decoding JSON: %w", err)} // err after errorf will be like "Error decoding JSON: 'Failed unmarshal: ...'
			return
		}
		result <- Response{str: fmt.Sprintf("Comic number: %d. Comic title: \"%s\".", comic.Num, comic.Title)}
	}

}

func main() {
	maxWorkers := flag.Int("workers", 5, "max number of work")
	left := flag.Int("first", 1, "left boundary")
	right := flag.Int("last", 100, "right boundary")
	flag.Parse()
	totalCnt := *right - *left + 1
	comicNum := make(chan int, totalCnt)
	result := make(chan Response, totalCnt)

	for w := 1; w <= *maxWorkers; w++ {
		go downloadComic(comicNum, result)
	}

	for n := *left; n <= *right; n++ {
		comicNum <- n
	}
	close(comicNum)

	for i := 1; i <= totalCnt; i++ {
		res := <-result
		if res.err != nil {
			fmt.Printf("%s", res.err.Error())
		}
		fmt.Println(res.str)
	}
}
