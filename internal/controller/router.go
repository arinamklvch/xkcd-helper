package controller

import (
	"net/http"

	_ "github.com/arinamklvch/xkcd-helper/docs"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/time/rate"
)

type userLimiters map[string]*rate.Limiter

// собираем таблицу маршрутов для сервера
// “распределитель” HTTP-запросов
func NewRouter(service *usecase.Service) *http.ServeMux {
	mux := http.NewServeMux()

	limiters := make(userLimiters)

	// обработчик с функциями из service
	handler := NewHandler(service)
	mux.HandleFunc("GET /load-comics", rateLimitMiddleware(limiters, handler.loadComics))
	mux.HandleFunc("GET /search-comics", rateLimitMiddleware(limiters, handler.searchComics))
	mux.HandleFunc("GET /last-comic", rateLimitMiddleware(limiters, handler.getLastComic))
	mux.HandleFunc("PUT /update", rateLimitMiddleware(limiters, handler.updateComics))
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}
