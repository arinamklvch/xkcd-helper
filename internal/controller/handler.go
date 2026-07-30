package controller

import "github.com/arinamklvch/xkcd-helper/internal/usecase"

// handler loads XKCD comics in the requested numeric range
type Handler struct {
	service *usecase.Service
}

// чтобы Handler мог пользоваться бизнес-логикой через service
func NewHandler(service *usecase.Service) *Handler {
	return &Handler{
		service: service,
	}
}
