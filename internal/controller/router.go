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
func NewRouter(service *usecase.Service, limit rate.Limit, burst int, JWTsecretKey string) *http.ServeMux {
	mux := http.NewServeMux()

	limiters := make(userLimiters)

	// обработчик с функциями из service
	handler := NewHandler(service)

	rateLimitMiddleware := func(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return rateLimit(limiters, limit, burst, next)
	}

	authMiddleware := func(needAdminCheck bool, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return auth(JWTsecretKey, needAdminCheck, next)
	}

	webAuthMiddleware := func(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return webAuth(JWTsecretKey, next)
	}

	// educational
	mux.HandleFunc("GET /load-comics", rateLimitMiddleware(authMiddleware(false, handler.loadComics)))
	mux.HandleFunc("GET /last-comic", rateLimitMiddleware(authMiddleware(false, handler.getLastComic)))
	mux.HandleFunc("PUT /update", rateLimitMiddleware(authMiddleware(true, handler.updateComics)))
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// useful
	mux.HandleFunc("GET /", rateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/login", http.StatusSeeOther) }))
	mux.HandleFunc("GET /login", rateLimitMiddleware(handler.webLogin))
	mux.HandleFunc("POST /login", rateLimitMiddleware(handler.login))
	mux.HandleFunc("GET /search-comics", rateLimitMiddleware(webAuthMiddleware(handler.searchComics)))

	return mux
}
