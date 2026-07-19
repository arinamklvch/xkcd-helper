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

func main() {
	left := flag.Int("first", 1, "left boundary")
	right := flag.Int("last", 100, "right boundary")
	flag.Parse()
	for *left <= *right {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", *left)
		response, err1 := http.Get(url)
		if err1 != nil {
			fmt.Println("Request error:", err1)
			continue
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusNotFound && *left != 404 {
			return
		}
		if response.StatusCode != http.StatusOK {
			fmt.Println("Error:", response.Status)
			*left++
			continue
		}
		body, err1 := io.ReadAll(response.Body)
		if err1 != nil {
			fmt.Println("Error reading response:", err1)
			*left++
			continue
		}
		var comic Comic
		err2 := json.Unmarshal(body, &comic)
		if err2 != nil {
			fmt.Println("Error decoding JSON:", err2)
			*left++
			continue
		}
		fmt.Printf("Comic number: %d. Comic title: \"%s\".\n", comic.Num, comic.Title)
		*left++
	}
}
