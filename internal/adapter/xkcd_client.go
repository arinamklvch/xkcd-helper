package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/arinamklvch/xkcd-helper/internal/domain"
)

const lastComicUrl = "https://xkcd.com/info.0.json"

type XkcdClient struct {
	client     http.Client
	maxWorkers int
	logger     *slog.Logger
}

func NewXkcdClient(client http.Client, maxWorkers int, logger *slog.Logger) *XkcdClient {
	return &XkcdClient{
		client:     client,
		maxWorkers: maxWorkers,
		logger:     logger,
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

func (x *XkcdClient) closeBody(body io.Closer) {
	if err := body.Close(); err != nil {
		x.logger.Error("failed to close response body", "error", err)
	}
}

func (x *XkcdClient) GetLastComicNum() (int, error) {
	resp, err := x.client.Get(lastComicUrl)
	if err != nil {
		x.logger.Error("failed to get last xkcd comic", "url", lastComicUrl, "error", err)
		return 0, err
	}
	defer x.closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %s", resp.Status)
		x.logger.Error("failed to get last xkcd comic", "url", lastComicUrl, "status", resp.Status, "error", err)
		return 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		x.logger.Error("failed to read last xkcd comic response body", "url", lastComicUrl, "error", err)
		return 0, err
	}

	var comic domain.Comic
	err = json.Unmarshal(body, &comic)
	if err != nil {
		x.logger.Error("failed to unmarshal last xkcd comic", "url", lastComicUrl, "error", err)
		return 0, err
	}

	return comic.Num, nil
}

func (x *XkcdClient) DownloadComicsRange(from, to int) ([]domain.Comic, error) {
	x.logger.Info("start downloading comics", "from", from, "to", to)
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
			x.logger.Error("failed to download comics range", "from", from, "to", to, "downloaded", count-1, "error", resp.Err)
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
			x.logger.Info("comics downloaded", "count", count, "total", totalCnt)
		}
	}
	x.logger.Info("finished downloading comics", "count", len(comics), "requested", totalCnt)
	return comics, nil
}

// worker идет в xkcd и скачивает комикс, ответ возвращает в канал responses,
// номер комикса принимает в comicNum
func (x *XkcdClient) worker(comicNum <-chan int, responses chan<- response) {
	for n := range comicNum {
		url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", n)
		resp, err := x.client.Get(url)

		if err != nil {
			x.logger.Error("failed to get xkcd comic", "comic_num", n, "url", url, "error", err)
			responses <- response{Err: err}
			continue
		}

		if n == http.StatusNotFound {
			x.closeBody(resp.Body)
			responses <- response{}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("unexpected status code: %s", resp.Status)
			x.logger.Error("failed to get xkcd comic", "comic_num", n, "url", url, "status", resp.Status, "error", err)
			x.closeBody(resp.Body)
			responses <- response{Err: err}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		x.closeBody(resp.Body)
		if err != nil {
			x.logger.Error("failed to read xkcd comic response body", "comic_num", n, "url", url, "error", err)
			responses <- response{Err: err}
			continue
		}

		var comic Comic
		err = json.Unmarshal(body, &comic)
		if err != nil {
			x.logger.Error("failed to unmarshal xkcd comic", "comic_num", n, "url", url, "error", err)
			responses <- response{Err: err}
			continue
		}

		responses <- response{Comic: &comic}
	}
}
