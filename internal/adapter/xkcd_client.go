package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
)

const latestComicUrl = "https://xkcd.com/info.0.json"
const maxWorkers = 5

type XkcdClient struct {
	client http.Client
}

func NewXkcdClient(client http.Client) *XkcdClient {
	return &XkcdClient{
		client: client,
	}
}

type Comic struct {
	Month      string
	Num        int
	Link       string
	Year       string
	News       string
	SafeTitle  string
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

type response struct {
	Comic *Comic
	Err   error
}

func (x *XkcdClient) GetLatestComicNum() (int, error) {
	resp, err := x.client.Get(latestComicUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Unexpected status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var comic domain.Comic
	err = json.Unmarshal(body, &comic)
	if err != nil {
		return 0, err
	}

	return comic.Num, nil
}

func (x *XkcdClient) DownloadComicsRange(from, to int) ([]domain.Comic, error) {
	totalCnt := to - from + 1
	comicNums := make(chan int, totalCnt)
	responses := make(chan response, totalCnt)

	for range maxWorkers {
		go x.worker(comicNums, responses)
	}

	for n := from; n <= to; n++ {
		comicNums <- n
	}
	close(comicNums)

	comics := make([]domain.Comic, 0, totalCnt)
	for range totalCnt {
		resp := <-responses
		if resp.Err != nil {
			return nil, resp.Err
		}
		if resp.Comic == nil {
			continue
		}

		comics = append(comics, domain.Comic{
			Month:      resp.Comic.Month,
			Num:        resp.Comic.Num,
			Link:       resp.Comic.Link,
			Year:       resp.Comic.Year,
			News:       resp.Comic.News,
			SafeTitle:  resp.Comic.SafeTitle,
			Transcript: resp.Comic.Transcript,
			Alt:        resp.Comic.Alt,
			Img:        resp.Comic.Img,
			Title:      resp.Comic.Title,
			Day:        resp.Comic.Day,
		})
	}

	return comics, nil
}

// worker идет в xkcd и скачивает комикс, ответ возвращает в канал responses,
// номер комикса принимает в comicNum
func (x *XkcdClient) worker(comicNum <-chan int, responses chan<- response) {
	for n := range comicNum {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", n)
		resp, err := x.client.Get(url)

		if err != nil {
			responses <- response{Err: err}
			continue
		}
		defer resp.Body.Close()

		if n == http.StatusNotFound {
			responses <- response{}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			responses <- response{Err: fmt.Errorf("Unexpected status code: %s", resp.Status)}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			responses <- response{Err: err}
			continue
		}

		var comic Comic
		err = json.Unmarshal(body, &comic)
		if err != nil {
			responses <- response{Err: err}
			continue
		}

		responses <- response{Comic: &comic}
	}
}
