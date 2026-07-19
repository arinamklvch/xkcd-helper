package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"safe_title"`
}

func main() {
	counter := 1
	for {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", counter)
		response, err1 := http.Get(url)
		if err1 != nil {
			fmt.Println("Request error:", err1)
			continue
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusNotFound && counter != 404 {
			return
		}
		if response.StatusCode != http.StatusOK {
			fmt.Println("Error:", response.Status)
			counter++
			continue
		}
		body, err1 := io.ReadAll(response.Body)
		if err1 != nil {
			fmt.Println("Error reading response:", err1)
			counter++
			continue
		}
		var comic Comic
		err2 := json.Unmarshal(body, &comic)
		if err2 != nil {
			fmt.Println("Error decoding JSON:", err2)
			counter++
			continue
		}
		fmt.Printf("Comic number: %d. Comic title: \"%s\".\n", comic.Num, comic.Title)
		counter++
	}
}
