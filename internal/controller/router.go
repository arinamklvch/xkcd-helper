package controller

import (
	"net/http"

	_ "github.com/arinamklvch/xkcd-helper/docs"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

// собираем таблицу маршрутов для сервера
// “распределитель” HTTP-запросов
func NewRouter(service *usecase.Service) *http.ServeMux {
	mux := http.NewServeMux()

	// обработчик с функциями из service
	handler := NewHandler(service)
	mux.HandleFunc("GET /load-comics", handler.loadComics)
	mux.HandleFunc("GET /search-comics", handler.searchComics)
	mux.HandleFunc("GET /last-comic", handler.getLastComic)
	mux.HandleFunc("PUT /update", handler.updateComics)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}
