package controller

import (
	"log/slog"

	"github.com/arinamklvch/xkcd-helper/internal/usecase"
)

// handler loads XKCD comics in the requested numeric range
type Handler struct {
	service *usecase.Service
	logger  *slog.Logger
}

// чтобы Handler мог пользоваться бизнес-логикой через service
func NewHandler(service *usecase.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
