package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
)

const lastComicUrl = "https://xkcd.com/info.0.json"

type XkcdClient struct {
	client     http.Client
	maxWorkers int
	// logger
}

func NewXkcdClient(client http.Client, maxWorkers int) *XkcdClient {
	return &XkcdClient{
		client:     client,
		maxWorkers: maxWorkers,
	}
}

type Comic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

type response struct {
	Comic *Comic
	Err   error
}

func closeBody(body io.Closer) {
	if err := body.Close(); err != nil {
		fmt.Println("failed to close response body:", err)
	}
}

func (x *XkcdClient) GetLastComicNum() (int, error) {
	resp, err := x.client.Get(lastComicUrl)
	if err != nil {
		return 0, err
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %s", resp.Status)
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
	fmt.Println("start downloading...")
	totalCnt := to - from + 1
	comicNums := make(chan int, totalCnt)
	responses := make(chan response, totalCnt)

	for range x.maxWorkers {
		go x.worker(comicNums, responses)
	}

	for n := from; n <= to; n++ {
		comicNums <- n
	}
	close(comicNums)

	comics := make([]domain.Comic, 0, totalCnt)
	var count int
	for range totalCnt {
		count++
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

		if count%100 == 0 {
			fmt.Println("comics downloaded:", count)
		}
	}
	fmt.Println("finished downloading.")
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

		if n == http.StatusNotFound {
			closeBody(resp.Body)
			responses <- response{}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			closeBody(resp.Body)
			responses <- response{Err: fmt.Errorf("unexpected status code: %s", resp.Status)}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		closeBody(resp.Body)
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
