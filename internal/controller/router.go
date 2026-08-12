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

	mux.HandleFunc("GET /load-comics", rateLimitMiddleware(limiters, authMiddleware(false, handler.loadComics)))
	mux.HandleFunc("GET /search-comics", rateLimitMiddleware(limiters, authMiddleware(false, handler.searchComics)))
	mux.HandleFunc("GET /last-comic", rateLimitMiddleware(limiters, authMiddleware(false, handler.getLastComic)))
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("PUT /update", rateLimitMiddleware(limiters, authMiddleware(true, handler.updateComics)))
	mux.HandleFunc("POST /login", rateLimitMiddleware(limiters, handler.login))
	return mux
}
