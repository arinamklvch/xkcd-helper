package controller

import (
	"net/http"

	_ "github.com/arinamklvch/xkcd-helper/docs"
	"github.com/arinamklvch/xkcd-helper/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(service *usecase.Service) *http.ServeMux {
	mux := http.NewServeMux()

	handler := NewHandler(service)
	mux.HandleFunc("GET /load-comics", handler.LoadComics)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}
