package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
)

type XkcdClient struct {
	client http.Client
}

func NewXkcdClient(client http.Client) *XkcdClient {
	return &XkcdClient{
		client: client,
	}
}

type response struct {
	Num   int
	Title string
	Err   error
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

// dto.Comic -- сущность для/из handler
func (x *XkcdClient) DownloadComics(from, to int) ([]domain.Comic, error) {
	totalCnt := from - to
	comicNums := make(chan int, totalCnt)
	responses := make(chan response, totalCnt)

	const maxWorkers = 5
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
			fmt.Println("Got error:", resp.Err)
			continue
		}
		comics = append(comics, domain.Comic{
			Num:   resp.Num,
			Title: resp.Title,
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
		responses <- response{Num: comic.Num, Title: comic.Title}
	}
}
