// Package main provides an HTTP API for loading XKCD comics by range.
//
// @title XKCD Helper API
// @version 1.0
// @description API for loading XKCD comic titles by comic number range.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/arinamklvch/xkcd-helper/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"safe_title"`
}

type Response struct {
	str string
	err error
}

func downloadComic(comicNum <-chan int, responses chan<- Response) {
	for n := range comicNum {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", n)
		response, err := http.Get(url)
		if err != nil {
			responses <- Response{err: err}
			continue
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			responses <- Response{err: fmt.Errorf("Unexpected status code: %s", response.Status)}
			continue
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			responses <- Response{err: err}
			continue
		}
		var comic Comic
		err = json.Unmarshal(body, &comic)
		if err != nil {
			responses <- Response{err: err}
			continue
		}
		responses <- Response{str: fmt.Sprintf("%d, %s", comic.Num, comic.Title)}
	}
}

// handler loads XKCD comics in the requested numeric range.
//
// @Summary Load comics
// @Description Returns XKCD comics as "number, title" strings for the requested range.
// @Tags comics
// @Produce json
// @Param from query int true "Starting comic number"
// @Param to query int true "Ending comic number"
// @Success 200 {array} string "Loaded comics"
// @Failure 500 {string} string "Failed to send JSON"
// @Router /load-comics [get]
func handler(w http.ResponseWriter, r *http.Request) {
	from, _ := strconv.Atoi(r.URL.Query().Get("from"))
	to, _ := strconv.Atoi(r.URL.Query().Get("to"))
	totalCnt := to - from + 1
	comicNums := make(chan int, totalCnt)
	responses := make(chan Response, totalCnt)
	maxWorkers := 5
	for range maxWorkers {
		go downloadComic(comicNums, responses)
	}
	for n := from; n <= to; n++ {
		comicNums <- n
	}
	close(comicNums)
	comics := []string{}
	for i := 1; i <= totalCnt; i++ {
		resp := <-responses
		if resp.err != nil {
			fmt.Println("Error:", resp.err)
			continue
		}
		comics = append(comics, resp.str)
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(comics)
	if err != nil {
		http.Error(w, "Failed to send JSON", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/load-comics", handler)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	server := http.Server{Addr: ":8081"}
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server failed:", err)
		}
	}()
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	s := <-signalChannel
	fmt.Println("\nCatched signal:", s)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
